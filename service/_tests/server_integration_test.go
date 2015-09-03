package integration_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
)

var (
	server   *negroni.Negroni
	request  *http.Request
	recorder *httptest.ResponseRecorder

	appEnv, _ = cfenv.Current()
)

var _ = Describe("Server Integration Test", func() {

	Describe("GET /v1/inventory", func() {
		BeforeEach(func() {
			server = NewServer(appEnv)
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest("GET", "/v1/inventory", nil)
		})

		Context("when inventory does not exist", func() {
			It("returns a status code of 200 and an empty set", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
				Ω(recorder.Body).To(ContainSubstring("[]"))
			})
		})

		// Context("when inventory exists", func() {
		// 	It("returns a status code of 200 and a list of inventory", func() {
		// 		server.ServeHTTP(recorder, request)
		// 		Ω(recorder.Code).To(Equal(200))
		// 		Ω(recorder.Body).To(ContainSubstring("2C.small"))
		// 	})
		// })
	})
})

//seed inventory collection
// i := &InventoryItem{}
// i.SKU = "2C.small"
// i.Tier = 2
// i.OfferingType = "C"
// i.Size = "small"
// i.Status = "available"
// b, _ := json.Marshal(i)
// reader := strings.NewReader(string(b))
// _, _ = http.NewRequest("POST", "/v1/inventory", reader)
