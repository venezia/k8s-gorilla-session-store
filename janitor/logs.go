package janitor

const (
	// Main Kubernetes Driver Messages
	LogCannotReadKubeconfigFile        = "cannot read in kubeconfig file as specified"
	LogCannotConvertKubeconfigContents = "could not convert kubeconfig contents to a valid configuration"
	LogCannotUseInClusterConfig        = "seemingly not in a kubernetes cluster, cannot use in cluster configuration"
	LogCannotConvertInClusterConfig    = "seemingly in a kubernetes cluster but cannot create a valid cluster configuration"

	LogSleepDurationTooShort = "sleep duration is set too short"

	LogK8SClientAlreadyExists                     = "kubernetes client already set, not going to recreate"
	LogCannotCreateKubernetesClient               = "error occurred while trying to create a kubernetes client"
	LogK8SConfigCreationThroughKubeconfigLocation = "attempting to create a kubernetes configuration via a kubeconfig file location setting"
	LogK8SConfigCreationThroughKubeconfig         = "attempting to create a kubernetes configuration via a kubeconfig contents provided directly in configuration"
	LogK8SConfigCreationThroughInCluster          = "attempting to create a kubernetes configuration through the in-cluster option"

	LogLeavingCleanContextClosed = "clean function ending as the parent context has closed"

	LogGettingListOfSessions                          = "getting list of sessions"
	LogCannotGetListOfSessions                        = "cannot get list of sessions"
	LogFoundSessionToDelete                           = "found session to delete because of age"
	LogCannotDeleteNonExistentSessionFromKubernetes   = "could not delete session from kubernetes - already removed"
	LogErrorTryingToDeleteSessionResourceInKubernetes = "error occurred while trying to delete a session resource in kubernetes"

	LogSessionIDKey          = "sessionID"
	LogKubeconfigLocationKey = "kubeconfigLocation"
	LogKubeconfigSizeKey     = "kubeconfigSize"
	LogErrorKey              = "error"
	LogContinueKey           = "continue"
)
