package gorilla_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/venezia/k8s-gorilla-session-store/gorilla"
)

const (
	interfaceSessionName = "interface-session"
)

var _ = Describe("Gorilla Session Interface Compliance Tests", Ordered, Label("k8s"), func() {
	Context("Ensure interface is working correctly", func() {
		store, err := gorilla.New(context.Background(), generateGenericConfig())
		Expect(err).ShouldNot(HaveOccurred())

		httpRequest, _ := http.NewRequest("GET", "/", nil)

		// New session functionality
		It("not return an error for a new session, and let me save it", func() {
			httpWriter := fakeHTTPResponseWriter{HeaderMap: make(http.Header)}

			session, err := store.New(httpRequest, interfaceSessionName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(session.IsNew).To(BeTrue())
			Expect(session.Values).To(BeEmpty())

			session.Values[sessionDataKey] = sessionDataValue
			err = store.Save(httpRequest, &httpWriter, session)
			Expect(err).ShouldNot(HaveOccurred())
			addCookie(httpRequest, &httpWriter)
		})
		It("should let me load a session through New, using the existing cookie", func() {
			session, err := store.New(httpRequest, interfaceSessionName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(session.IsNew).To(BeFalse())
			Expect(len(session.Values)).To(Equal(1))

			err = store.K8SDelete(context.Background(), session)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

type fakeHTTPResponseWriter struct {
	HeaderMap http.Header
}

func (f *fakeHTTPResponseWriter) Header() http.Header {
	return f.HeaderMap
}

func (f *fakeHTTPResponseWriter) Write(input []byte) (int, error) {
	return len(input), nil
}

func (f *fakeHTTPResponseWriter) WriteHeader(statusCode int) {}

func addCookie(r *http.Request, w *fakeHTTPResponseWriter) {
	cookieMap := w.HeaderMap.Get("Set-Cookie")
	if cookieMap != "" {
		header := http.Header{}
		header.Add("Cookie", cookieMap)
		dummyRequest := http.Request{Header: header}
		cookieList := dummyRequest.Cookies()
		for _, cookie := range cookieList {
			r.AddCookie(cookie)
		}
	}
}
