package gorilla_test

import (
	"context"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"golang.org/x/exp/slog"

	"github.com/venezia/k8s-gorilla-session-store/gorilla"
)

const (
	sessionDataKey   = "foo"
	sessionDataValue = "bario"
)

var _ = Describe("K8S Compliance Tests", Label("k8s"), func() {

	Describe("EnsureClient", func() {
		Context("Ensure the client is working", func() {
			It("not return an error", func() {
				_, err := gorilla.New(context.Background(), generateGenericConfig())

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("K8S CRUD", Ordered, func() {
		Context("Ensure CRUD is working correctly", func() {
			store, err := gorilla.New(context.Background(), generateGenericConfig())
			Expect(err).ShouldNot(HaveOccurred())

			// Save Functionality
			It("not return an error for a sane save", func() {
				testSession := sessions.Session{
					ID:      "k8s-save-test",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
					IsNew:   true,
				}
				testSession.Values[sessionDataKey] = sessionDataValue
				err = store.K8SSave(context.Background(), &testSession)
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should return an error for an illegal save", func() {
				testSession := sessions.Session{
					ID:      "too-many-characters-this-is-too-many-characters-this-is-too-many-character-this-is-too-many-characters-this-is-too-many-characters-this-is-too-many-characters-this-is-too-many-characters-this-is-too-many-characters-this-is-too-many-characters-this-is-too-many-characters-this-is-too-many-characters",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
					IsNew:   true,
				}
				testSession.Values["foo"] = "bario"
				err = store.K8SSave(context.Background(), &testSession)
				Expect(err).Should(HaveOccurred())
			})

			// Load Functionality
			It("should return a previously saved session", func() {
				testSession := sessions.Session{
					ID:      "k8s-save-test",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
				}
				fetchSucccess, err := store.K8SLoad(context.Background(), &testSession)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fetchSucccess).To(BeTrue())
				fooValue, ok := testSession.Values[sessionDataKey].(string)
				Expect(ok).To(BeTrue())
				Expect(fooValue).To(Equal(sessionDataValue))
			})
			It("should return an empty map for a missing session", func() {
				testSession := sessions.Session{
					ID:      "k8s-save-test-missing",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
				}
				fetchSucccess, err := store.K8SLoad(context.Background(), &testSession)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fetchSucccess).To(BeFalse())
				Expect(len(testSession.Values)).To(Equal(0))
			})
			It("should return an error for an invalid session name", func() {
				testSession := sessions.Session{
					ID:      "",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
				}
				fetchSucccess, err := store.K8SLoad(context.Background(), &testSession)
				Expect(err).Should(HaveOccurred())
				Expect(fetchSucccess).To(BeFalse())
				Expect(len(testSession.Values)).To(Equal(0))
			})

			// Delete Functionality
			It("not return an error for a sane delete", func() {
				testSession := sessions.Session{
					ID:      "k8s-save-test",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
					IsNew:   true,
				}
				testSession.Values["foo"] = "bario"
				err = store.K8SDelete(context.Background(), &testSession)
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should return an error for an illegal delete", func() {
				testSession := sessions.Session{
					ID:      "",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
					IsNew:   true,
				}
				err = store.K8SDelete(context.Background(), &testSession)
				Expect(err).Should(HaveOccurred())
			})
			It("should safely ignore a delete for something that never existed", func() {
				testSession := sessions.Session{
					ID:      "i-do-not-exist",
					Values:  make(map[interface{}]interface{}),
					Options: &sessions.Options{},
					IsNew:   true,
				}
				err = store.K8SDelete(context.Background(), &testSession)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})

func generateGenericConfig() *gorilla.Config {
	config := &gorilla.Config{
		K8S: gorilla.K8SConfig{
			KubeconfigLocation: os.Getenv(kubeconfigEnvironmentVariable),
			Namespace:          getKubernetesNamespace(),
		},
		SessionOptions: &sessions.Options{},
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
