package gorilla

const (
	// Main Kubernetes Driver Messages
	LogCannotReadKubeconfigFile        = "cannot read in kubeconfig file as specified"
	LogCannotConvertKubeconfigContents = "could not convert kubeconfig contents to a valid configuration"
	LogCannotUseInClusterConfig        = "seemingly not in a kubernetes cluster, cannot use in cluster configuration"
	LogCannotConvertInClusterConfig    = "seemingly in a kubernetes cluster but cannot create a valid cluster configuration"

	LogSavingSessionDataToKubernetes                  = "saving session data to kubernetes"
	LogCannotPersistDataToKubernetes                  = "could not save data to kubernetes"
	LogLoadSessionDataFromKubernetes                  = "loading session data from kubernetes"
	LogSessionDataNotFoundInKubernetes                = "session resource not found in kubernetes"
	LogSessionDataRetrievalErrorOccurred              = "error occurred while trying to retrieve session data from kubernetes"
	LogErrorTryingToDeserializeSessionData            = "error occurred while trying to deserialize session data"
	LogCannotDeleteNonExistentSessionFromKubernetes   = "could not delete session from kubernetes - already removed"
	LogDeletingSessionFromKubernetes                  = "deleting session from kubernetes"
	LogErrorTryingToDeleteSessionResourceInKubernetes = "error occurred while trying to delete a session resource in kubernetes"

	LogK8SClientAlreadyExists                     = "kubernetes client already set, not going to recreate"
	LogCannotCreateKubernetesClient               = "error occurred while trying to create a kubernetes client"
	LogK8SConfigCreationThroughKubeconfigLocation = "attempting to create a kubernetes configuration via a kubeconfig file location setting"
	LogK8SConfigCreationThroughKubeconfig         = "attempting to create a kubernetes configuration via a kubeconfig contents provided directly in configuration"
	LogK8SConfigCreationThroughInCluster          = "attempting to create a kubernetes configuration through the in-cluster option"

	LogSessionIDKey          = "sessionID"
	LogKubeconfigLocationKey = "kubeconfigLocation"
	LogKubeconfigSizeKey     = "kubeconfigSize"

	ErrSessionDataTooBig = "k8s-gorilla-session-store.K8SSave: data to persist is over size limit"
)
