/*
Package gorilla provides a [session.Store] implementation that uses kubernetes as its persistence layer

# Purpose

This package will create sessions resources inside a kubernetes cluster.  It does not clean the sessions up.
Please consider using the package [github.com/venezia/k8s-gorilla-session-store/janitor] to clean up expired sessions.

# Kubernetes Considerations

  - This package uses a generated client [github.com/venezia/k8s-gorilla-session-store/k8s/client/clientset/versioned.Interface] to interact with kubernetes.
  - There is a custom resource definition file provided in the k8s/crd folder.  This needs to be applied to the kubernetes cluster
  - If RBAC is enabled, RBAC permissions will probably be needed

# Logging Considerations

This package supports [log/slog.Logger] as a logger provided to it through the configuration of a new store.

# Size limitations

Generally speaking, kubernetes resources are limited to 1MiB in size.  Therefore, a default limit of 128KiB is used on
session storage.  One can set the limit as they wish, but the size may exceed what the kubernetes apiserver accepts.

[session.Store]: https://pkg.go.dev/github.com/gorilla/sessions#Store
*/
package gorilla
