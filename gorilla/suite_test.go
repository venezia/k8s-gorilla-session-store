package gorilla_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	kubeconfigEnvironmentVariable          = "KUBECONFIG"
	kubernetesNamespaceEnvironmentVariable = "K8S_SESSION_NAMESPACE"
	kubernetesNamespaceDefault             = "default"
	loggingEnabledEnvironmentVariable      = "LOGGING_ENABLED"
)

func getKubernetesNamespace() string {
	namespace := os.Getenv(kubernetesNamespaceEnvironmentVariable)
	if namespace != "" {
		return namespace
	}
	return kubernetesNamespaceDefault
}

func KubeEnabled() bool {
	return os.Getenv(kubeconfigEnvironmentVariable) != ""
}

func TestProcessor(t *testing.T) {
	_, reporterConfig := GinkgoConfiguration()
	RegisterFailHandler(Fail)

	if KubeEnabled() {
		RunSpecs(t, "Gorilla K8S Session Suite", reporterConfig, Label("gorilla"))
	}
}
