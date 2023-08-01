package gorilla

import (
	"encoding/base32"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Get should return a cached session.
func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New creates new session and if possible, loads the stored data into it
func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	var err error
	var ok bool

	// Create a new session and K8SLoad session options into it
	session := sessions.NewSession(s, name)
	session.Options = s.Config.SessionOptions
	// Mark the session as new for now
	session.IsNew = true

	// try to load a cookie, if so continue
	if c, errCookie := r.Cookie(name); errCookie == nil {
		// cookie was retrieved, let's try to decode the session id
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			// Let's try to load the session data from the k8s store
			ok, err = s.K8SLoad(r.Context(), session)
			// if ok is true and there is no error, we were successful
			if err == nil && ok {
				session.IsNew = false
			}
		}
	}
	return session, err

}

// Save should persist session to the store (kubernetes)
func (s *Store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// if max age is greater than zero, we want to persist the data to the store (kubernetes)
	if session.Options.MaxAge > 0 {
		// Check if an id has not yet been given
		if session.ID == "" {
			// Generates random id
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}

		// Let's try persisting the data to the store (kubernetes)
		if err := s.K8SSave(r.Context(), session); err != nil {
			return err
		}

		// Now let's generate the contents to persist as cookie to the browser
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	} else {
		// So we need to K8SDelete the data from k8s and browser
		// First let's delete the key from k8s...
		if err := s.K8SDelete(r.Context(), session); err != nil {
			return err
		}
		// Now let's remove the session name from the browser
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	}
	return nil
}
