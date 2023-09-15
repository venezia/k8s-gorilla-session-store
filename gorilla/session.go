package gorilla

import (
	"context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log/slog"

	clientset "github.com/venezia/k8s-gorilla-session-store/k8s/client/clientset/versioned"
)

const (
	DefaultMaxAge     = 86400             // how long should the session persist for
	DefaultMaxLength  = 1024 * 128        // 128KB of session storage
	DefaultCookiePath = "/"               // have the session be site wide
	DefaultSecret     = "notAGreatSecret" // default secret for session validation
)

type Config struct {
	K8S            K8SConfig
	SessionOptions *sessions.Options
	MaxLength      int    // Max amount of storage to be allowed in session data for a given user
	Secret         string // Secret is for validating a session key
	Logger         *slog.Logger
}

type K8SConfig struct {
	KubeconfigLocation string // Where a kubeconfig file is located on the system
	Kubeconfig         []byte // The contents of a kubeconfig
	Namespace          string // The namespace that will be used to look for resources
	Client             clientset.Interface
}

type Store struct {
	Ctx        context.Context
	Config     *Config
	Serializer SessionSerializer
	Codecs     []securecookie.Codec
}

// New returns a new Store based on configuration
func New(ctx context.Context, config *Config) (*Store, error) {
	store := &Store{
		Config:     config,
		Ctx:        ctx,
		Serializer: GobSerializer{},
	}
	applyDefaultSessionConfiguration(store.Config)
	store.Codecs = securecookie.CodecsFromPairs([]byte(config.Secret))

	err := store.EnsureK8SClient()
	if err != nil {
		return nil, err
	}

	return store, nil
}

// Apply reasonable defaults to the options
func applyDefaultSessionConfiguration(config *Config) {
	if config.MaxLength == 0 {
		config.MaxLength = DefaultMaxLength
	}
	if config.SessionOptions.Path == "" {
		config.SessionOptions.Path = DefaultCookiePath
	}
	if config.SessionOptions.MaxAge == 0 {
		config.SessionOptions.MaxAge = DefaultMaxAge
	}
	if config.Secret == "" {
		config.Secret = DefaultSecret
	}
}
