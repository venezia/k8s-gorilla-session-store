package gorilla

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"github.com/gorilla/sessions"
)

// SessionSerializer provides an interface hook for alternative serializers
type SessionSerializer interface {
	Deserialize(s string, session *sessions.Session) error
	Serialize(session *sessions.Session) (string, error)
}

// GobSerializer uses gob package to encode the session map
type GobSerializer struct{}

// Serialize using gob
func (g GobSerializer) Serialize(session *sessions.Session) (string, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(session.Values)
	if err == nil {
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	}
	return "", err
}

// Deserialize back to map[interface{}]interface{}
func (g GobSerializer) Deserialize(input string, session *sessions.Session) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(bytes.NewBuffer(decodedBytes))
	return dec.Decode(&session.Values)
}
