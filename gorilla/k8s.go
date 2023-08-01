package gorilla

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	v1alpha1apply "github.com/venezia/k8s-gorilla-session-store/k8s/client/applyconfiguration/session/v1alpha1"
	clientset "github.com/venezia/k8s-gorilla-session-store/k8s/client/clientset/versioned"
)

// K8SSave Stores the session in the kubernetes api.
func (s *Store) K8SSave(ctx context.Context, session *sessions.Session) error {
	serializedSession, err := s.Serializer.Serialize(session)
	if err != nil {
		return err
	}

	// Ensuring that we're not storing too much data
	if s.Config.MaxLength != 0 && len(serializedSession) > s.Config.MaxLength {
		return errors.New(ErrSessionDataTooBig)
	}

	// We need pointers to strings for the Apply call
	kind := "Session"
	apiVersion := "gorilla.michaelvenezia.com/v1alpha1"
	sessionID := NormalizeIDValue(session.ID)

	// Let's do an Apply call with the session data
	if s.Config.Logger != nil {
		s.Config.Logger.Debug(LogSavingSessionDataToKubernetes, LogSessionIDKey, sessionID)
	}
	_, err = s.Config.K8S.Client.GorillaV1alpha1().Sessions(s.Config.K8S.Namespace).Apply(ctx, &v1alpha1apply.SessionApplyConfiguration{
		TypeMetaApplyConfiguration: applymetav1.TypeMetaApplyConfiguration{
			Kind:       &kind,
			APIVersion: &apiVersion,
		},
		ObjectMetaApplyConfiguration: &applymetav1.ObjectMetaApplyConfiguration{
			Name: &sessionID,
		},
		Spec: &v1alpha1apply.SessionSpecApplyConfiguration{
			Data: &serializedSession,
		},
		Status: &v1alpha1apply.SessionStatusApplyConfiguration{
			TTL: &metav1.Time{time.Now().Add(time.Duration(s.Config.SessionOptions.MaxAge) * time.Second)},
		},
	}, metav1.ApplyOptions{FieldManager: "session-manager", Force: true})

	// Check for errors and wrap them
	if err != nil {
		return &Error{
			Type:    ExecutionErrorType,
			Message: LogCannotPersistDataToKubernetes,
			Err:     err,
		}
	}

	return nil
}

// K8SLoad reads the session from the kubernetes api.
// returns true if there is session data to be found
func (s *Store) K8SLoad(ctx context.Context, session *sessions.Session) (bool, error) {
	if s.Config.Logger != nil {
		s.Config.Logger.Debug(LogLoadSessionDataFromKubernetes, LogSessionIDKey, NormalizeIDValue(session.ID))
	}
	record, err := s.Config.K8S.Client.GorillaV1alpha1().Sessions(s.Config.K8S.Namespace).Get(ctx, NormalizeIDValue(session.ID), metav1.GetOptions{})
	if err != nil {
		// Not Found is OK - we may not have the session data for a variety of reasons
		if k8serrors.IsNotFound(err) {
			if s.Config.Logger != nil {
				s.Config.Logger.Debug(LogSessionDataNotFoundInKubernetes, LogSessionIDKey, NormalizeIDValue(session.ID))
			}
			return false, nil
		}
		// It wasn't not found, so we need to return an error...
		return false, &Error{
			Type:    ExecutionErrorType,
			Message: LogSessionDataRetrievalErrorOccurred,
			Err:     err,
		}
	}

	// If there is no data to deserialize, don't both deserializing it...
	if record.Spec.Data == "" {
		return false, nil
	}

	// Return a positive hit, and deserialize the data, returning any error
	deserializeError := s.Serializer.Deserialize(record.Spec.Data, session)
	if deserializeError != nil {
		return true, &Error{
			Type:    ExecutionErrorType,
			Message: LogErrorTryingToDeserializeSessionData,
			Err:     deserializeError,
		}
	}
	return true, nil
}

// K8SDelete removes a session's contents from the data store (kubernetes)
func (s *Store) K8SDelete(ctx context.Context, session *sessions.Session) error {
	if s.Config.Logger != nil {
		s.Config.Logger.Debug(LogDeletingSessionFromKubernetes, LogSessionIDKey, NormalizeIDValue(session.ID))
	}
	err := s.Config.K8S.Client.GorillaV1alpha1().Sessions(s.Config.K8S.Namespace).Delete(ctx, NormalizeIDValue(session.ID), metav1.DeleteOptions{})

	// No errors, move on
	if err == nil {
		return nil
	}

	// Error occured, but let's see if it was just deleting something that did not exist...
	if k8serrors.IsNotFound(err) {
		if s.Config.Logger != nil {
			s.Config.Logger.Debug(LogCannotDeleteNonExistentSessionFromKubernetes, LogSessionIDKey, NormalizeIDValue(session.ID))
		}
		return nil
	}

	// A real error occurred, wrap it and return it
	return &Error{
		Type:    ExecutionErrorType,
		Message: LogErrorTryingToDeleteSessionResourceInKubernetes,
		Err:     err,
	}
}

// EnsureK8SClient will try to create a kubernetes client if it isn't already set
func (s *Store) EnsureK8SClient() error {
	if s.Config.K8S.Client != nil {
		if s.Config.Logger != nil {
			s.Config.Logger.Debug(LogK8SClientAlreadyExists)
		}
		return nil
	}

	// OK, no client is set, so let's set it for us...
	config, err := s.GetK8SConfig()
	if err != nil {
		return err
	}

	s.Config.K8S.Client, err = clientset.NewForConfig(config)
	if err != nil {
		return &Error{
			Type:    ExecutionErrorType,
			Message: LogCannotCreateKubernetesClient,
			Err:     err,
		}
	}
	return nil
}

// GetK8SConfig will try to get a working kubernetes client configuration returning an error if it cannot
func (s *Store) GetK8SConfig() (*rest.Config, error) {
	var err error
	var config *rest.Config

	// Let's see if an explicit location for kubeconfig has been provided...
	if s.Config.K8S.KubeconfigLocation != "" {
		if s.Config.Logger != nil {
			s.Config.Logger.Debug(LogK8SConfigCreationThroughKubeconfigLocation, LogKubeconfigLocationKey, s.Config.K8S.KubeconfigLocation)
		}
		s.Config.K8S.Kubeconfig, err = os.ReadFile(s.Config.K8S.KubeconfigLocation)
		if err != nil {
			return nil, &Error{
				Type:    InputErrorType,
				Message: LogCannotReadKubeconfigFile,
				Err:     err,
			}
		}
	}

	// Let's see if we now have the raw contents of a kubeconfig...
	if len(s.Config.K8S.Kubeconfig) > 0 {
		if s.Config.Logger != nil {
			s.Config.Logger.Debug(LogK8SConfigCreationThroughKubeconfig, LogKubeconfigSizeKey, len(s.Config.K8S.Kubeconfig))
		}
		config, err = clientcmd.RESTConfigFromKubeConfig(s.Config.K8S.Kubeconfig)
		if err != nil {
			return nil, &Error{
				Type:    InputErrorType,
				Message: LogCannotConvertKubeconfigContents,
				Err:     err,
			}
		}
	} else {
		// No explicit kubeconfig was specified, let's see if we're in a cluster already...
		if s.Config.Logger != nil {
			s.Config.Logger.Debug(LogK8SConfigCreationThroughInCluster)
		}
		config, err = rest.InClusterConfig()
		if errors.Is(err, rest.ErrNotInCluster) {
			return nil, &Error{
				Type:    InputErrorType,
				Message: LogCannotUseInClusterConfig,
				Err:     err,
			}
		}
		if err != nil {
			return nil, &Error{
				Type:    ExecutionErrorType,
				Message: LogCannotConvertInClusterConfig,
				Err:     err,
			}
		}
	}

	return config, nil
}

// NormalizeIDValue will return back a lowercase string.  Kubernetes id values have to be lowercase
func NormalizeIDValue(input string) string {
	return strings.ToLower(input)
}
