/*
Package janitor will clean up expired sessions

# Purpose

Generally speaking, this package will remove presumably expired session resources from a kubernetes cluster.  It is intended
to compliment the [github.com/venezia/k8s-gorilla-session-store/gorilla] package.  This package could be run by the same
service that is creating the sessions, or simply be a cron-job or supporting service that is just cleaning up the detritus

# Intended Use

The package is meant to be used like so:

	    ctx, ctxCancel := context.WithCancel(context.Background())
		janitor, err := janitor.New(ctx, config)
	    go func() { janitor.Clean() }()
	    go func() {
	        time.Sleep(120 * time.Second)
	        ctxCancel()
	    }()

Where the janitor would run for 120 seconds.  Obviously one may want it to run forever, etc.

# Kubernetes Considerations

  - This package uses a generated client [github.com/venezia/k8s-gorilla-session-store/k8s/client/clientset/versioned.Interface] to interact with kubernetes.
  - There is a custom resource definition file provided in the k8s/crd folder.  This needs to be applied to the kubernetes cluster
  - If RBAC is enabled, RBAC permissions will probably be needed

# Errors

While an error can be generated if a Janitor is misconfigured, errors are not going to be returned to the caller.
Perhaps at a later time how to handle errors for a background job could be fleshed out.  Instead when something that
would typically result in an error being returned happens, a log entry is provided.

# Logging Considerations

This package supports [golang.org/x/exp/slog.Logger] as a logger provided to it through the configuration of a new store.
This will transition into the standard library implementation once go 1.21 is released
*/
package janitor
