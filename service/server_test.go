package pezinventory_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/codegangsta/negroni"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
)

type mockHandler struct{}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	next(w, req)
}

var _ = Describe("Server", func() {
	var (
		server   *negroni.Negroni
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		//handler := &mockHandler{}
		// server = NewServer(handler)
		server = NewServer()
		recorder = httptest.NewRecorder()
	})

	Describe("GET /inventory", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/inventory", nil)
		})

		// Context("when inventory does not exist", func() {
		// 	It("returns a status code of 200 and an empty set", func() {
		// 		server.ServeHTTP(recorder, request)
		// 		Ω(recorder.Code).To(Equal(200))
		// 		Ω(recorder.Body).To(ContainSubstring("[]"))
		// 	})
		// })

		Context("when inventory exists", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
			})

			It("returns a list of types", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Body).To(ContainSubstring("2C.small"))
			})
		})
	})
})
