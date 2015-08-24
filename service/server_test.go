package pezinventory_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-martini/martini"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
)

var _ = Describe("Server", func() {
	var (
		server   *martini.ClassicMartini
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		var f martini.Handler = func() {}
		server = NewServer(f)
		recorder = httptest.NewRecorder()
	})

	Describe("GET /v1/types", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/types", nil)
		})

		Context("when types exist", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
			})

			It("returns a list of types", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Body).To(ContainSubstring("success"))
				Ω(recorder.Body).To(ContainSubstring("data"))
			})
		})
	})

	Describe("GET /v1/types/:id", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/types/12345", nil)
		})

		Context("with valid Type id", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
			})

			It("returns Type details", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Body).To(ContainSubstring("success"))
				Ω(recorder.Body).To(ContainSubstring("data"))
			})
		})
	})

	Describe("GET /v1/types/:id/items", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/types/12345/items", nil)
		})

		Context("when Items exist for Type", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
			})

			It("returns a list of Items for Type", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Body).To(ContainSubstring("success"))
				Ω(recorder.Body).To(ContainSubstring("data"))
			})
		})
	})

	Describe("GET /v1/items", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/items", nil)
		})

		Context("when Items exist", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
			})

			It("returns a list of Items", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Body).To(ContainSubstring("success"))
				Ω(recorder.Body).To(ContainSubstring("data"))
			})
		})
	})

	Describe("GET /v1/items/:id", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/items/12345", nil)
		})

		Context("with valid Item id", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
			})

			It("returns Item details", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Body).To(ContainSubstring("success"))
				Ω(recorder.Body).To(ContainSubstring("data"))
			})
		})
	})

	Describe("GET /v1/items/:id/history", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/v1/items/12345/history", nil)
		})

		Context("not implemented", func() {
			It("returns a status code of 501", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(501))
			})

			It("returns error response", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Body).To(ContainSubstring("error"))
				Ω(recorder.Body).To(ContainSubstring("message"))
			})
		})
	})
})
