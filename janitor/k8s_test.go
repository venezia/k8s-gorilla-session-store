package janitor_test

import (
	"context"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"golang.org/x/exp/slog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"

	"github.com/venezia/k8s-gorilla-session-store/janitor"
	v1alpha1apply "github.com/venezia/k8s-gorilla-session-store/k8s/client/applyconfiguration/session/v1alpha1"
	"github.com/venezia/k8s-gorilla-session-store/k8s/client/clientset/versioned"
)

var _ = Describe("K8S Compliance Tests", Label("k8s"), func() {

	Describe("EnsureClient", func() {
		Context("Ensure the client is working", func() {
			It("not return an error", func() {
				_, err := janitor.New(context.Background(), generateGenericConfig())

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Janitor Cleanup", Ordered, func() {
		Context("Ensure cleanup is working correctly", func() {
			janitor, err := janitor.New(context.Background(), generateGenericConfig())
			Expect(err).ShouldNot(HaveOccurred())
			err = createOldSession(janitor.Ctx, janitor.Config.K8S.Client, janitor.Config.K8S.Namespace)
			It("not return an error for a sane save", func() {
				janitor.K8SCleanExpiredSessions()
			})
		})
	})
})

func generateGenericConfig() *janitor.Config {
	config := &janitor.Config{
		K8S: janitor.K8SConfig{
			KubeconfigLocation: os.Getenv(kubeconfigEnvironmentVariable),
			Namespace:          getKubernetesNamespace(),
		},
		SleepDuration: 3 * time.Second,
	}

	switch strings.ToLower(os.Getenv(loggingEnabledEnvironmentVariable)) {
	case "true", "1":
		logConfig := zap.NewDevelopmentConfig()
		logConfig.OutputPaths = []string{"stdout"}
		logConfig.ErrorOutputPaths = []string{"stderr"}
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logConfig.Encoding = "console"
		zapLogger, _ := logConfig.Build()

		config.Logger = slog.New(zapslog.NewHandler(zapLogger.Core()))
	}

	return config
}

func createOldSession(ctx context.Context, client versioned.Interface, namespace string) error {
	kind := "Session"
	apiVersion := "gorilla.michaelvenezia.com/v1alpha1"
	sessionID := "test-old-session"

	_, err := client.GorillaV1alpha1().Sessions(namespace).Apply(ctx, &v1alpha1apply.SessionApplyConfiguration{
		TypeMetaApplyConfiguration: applymetav1.TypeMetaApplyConfiguration{
			Kind:       &kind,
			APIVersion: &apiVersion,
		},
		ObjectMetaApplyConfiguration: &applymetav1.ObjectMetaApplyConfiguration{
			Name: &sessionID,
		},
		Status: &v1alpha1apply.SessionStatusApplyConfiguration{
			TTL: &metav1.Time{Time: time.Now().Add(-120 * time.Second)},
		},
	}, metav1.ApplyOptions{FieldManager: "session-manager", Force: true})
	return err
}
