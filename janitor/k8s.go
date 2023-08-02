package janitor

import (
	"errors"
	"os"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/venezia/k8s-gorilla-session-store/k8s/client/clientset/versioned"
)

// K8SCleanExpiredSessions will fetch all sessions, look to see if any of them have expired, and then remove them.
// The list is paginated - currently set through constant SessionListPaginationSize - but it will iterate across all sessions
func (j *Janitor) K8SCleanExpiredSessions() {
	continueValue := ""
	for {
		if j.Config.Logger != nil {
			j.Config.Logger.Debug(LogGettingListOfSessions, LogContinueKey, continueValue)
		}
		list, err := j.Config.K8S.Client.GorillaV1alpha1().Sessions(j.Config.K8S.Namespace).List(j.Ctx, metav1.ListOptions{Limit: SessionListPaginationSize, Continue: continueValue})
		if err != nil && j.Config.Logger != nil {
			j.Config.Logger.Warn(LogCannotGetListOfSessions, LogErrorKey, err)
			return
		}
		continueValue = list.Continue
		for _, item := range list.Items {
			if item.Status.TTL.Time.Before(time.Now()) {
				if j.Config.Logger != nil {
					j.Config.Logger.Debug(LogFoundSessionToDelete, LogSessionIDKey, item.Name)
				}
				err = j.Config.K8S.Client.GorillaV1alpha1().Sessions(j.Config.K8S.Namespace).Delete(j.Ctx, item.Name, metav1.DeleteOptions{})

				// Check for errors
				if err != nil && j.Config.Logger != nil {
					// Error occurred, but let's see if it was just deleting something that did not exist...
					if k8serrors.IsNotFound(err) {
						j.Config.Logger.Debug(LogCannotDeleteNonExistentSessionFromKubernetes, LogSessionIDKey, item.Name)
					} else {
						j.Config.Logger.Warn(LogErrorTryingToDeleteSessionResourceInKubernetes, LogErrorKey, err)

					}
				}
			}
		}
		if continueValue == "" {
			break
		}
	}
}

// EnsureK8SClient will try to create a kubernetes client if it isn't already set
func (j *Janitor) EnsureK8SClient() error {
	if j.Config.K8S.Client != nil {
		if j.Config.Logger != nil {
			j.Config.Logger.Debug(LogK8SClientAlreadyExists)
		}
		return nil
	}

	// OK, no client is set, so let's set it for us...
	config, err := j.GetK8SConfig()
	if err != nil {
		return err
	}

	j.Config.K8S.Client, err = clientset.NewForConfig(config)
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
func (j *Janitor) GetK8SConfig() (*rest.Config, error) {
	var err error
	var config *rest.Config

	// Let's see if an explicit location for kubeconfig has been provided...
	if j.Config.K8S.KubeconfigLocation != "" {
		if j.Config.Logger != nil {
			j.Config.Logger.Debug(LogK8SConfigCreationThroughKubeconfigLocation, LogKubeconfigLocationKey, j.Config.K8S.KubeconfigLocation)
		}
		j.Config.K8S.Kubeconfig, err = os.ReadFile(j.Config.K8S.KubeconfigLocation)
		if err != nil {
			return nil, &Error{
				Type:    InputErrorType,
				Message: LogCannotReadKubeconfigFile,
				Err:     err,
			}
		}
	}

	// Let's see if we now have the raw contents of a kubeconfig...
	if len(j.Config.K8S.Kubeconfig) > 0 {
		if j.Config.Logger != nil {
			j.Config.Logger.Debug(LogK8SConfigCreationThroughKubeconfig, LogKubeconfigSizeKey, len(j.Config.K8S.Kubeconfig))
		}
		config, err = clientcmd.RESTConfigFromKubeConfig(j.Config.K8S.Kubeconfig)
		if err != nil {
			return nil, &Error{
				Type:    InputErrorType,
				Message: LogCannotConvertKubeconfigContents,
				Err:     err,
			}
		}
	} else {
		// No explicit kubeconfig was specified, let's see if we're in a cluster already...
		if j.Config.Logger != nil {
			j.Config.Logger.Debug(LogK8SConfigCreationThroughInCluster)
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
