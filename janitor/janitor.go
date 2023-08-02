package janitor

import (
	"context"
	"time"

	"golang.org/x/exp/slog"

	clientset "github.com/venezia/k8s-gorilla-session-store/k8s/client/clientset/versioned"
)

const (
	MinimumSleepDuration      = 1 * time.Second
	SessionListPaginationSize = 50
)

type Config struct {
	K8S           K8SConfig
	SleepDuration time.Duration
	Logger        *slog.Logger
}

type K8SConfig struct {
	KubeconfigLocation string // Where a kubeconfig file is located on the system
	Kubeconfig         []byte // The contents of a kubeconfig
	Namespace          string // The namespace that will be used to look for resources
	Client             clientset.Interface
}

type Janitor struct {
	Ctx    context.Context
	Config *Config
}

// New returns a new Janitor based on configuration
func New(ctx context.Context, config *Config) (*Janitor, error) {
	// Sanity check...
	if config.SleepDuration < MinimumSleepDuration {
		return nil, &Error{
			Type:    InputErrorType,
			Message: LogSleepDurationTooShort,
		}
	}

	janitor := &Janitor{
		Config: config,
		Ctx:    ctx,
	}

	err := janitor.EnsureK8SClient()
	if err != nil {
		return nil, err
	}

	return janitor, nil
}

func (j *Janitor) Clean() {
	select {
	case <-j.Ctx.Done():
		if j.Config.Logger != nil {
			j.Config.Logger.Info(LogLeavingCleanContextClosed)
		}
	default:
		j.K8SCleanExpiredSessions()
		time.Sleep(j.Config.SleepDuration)
	}
}
