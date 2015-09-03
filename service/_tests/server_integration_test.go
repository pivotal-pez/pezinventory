package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
)

const (
	inventoryDB      = "inventory-db"
	inventoryURI     = "uri"
	fakeURIFormatter = "mongodb://%s:%s/25210aaf-xxxx-xxxx-9d53-0f4d078b1768"

	//VcapServicesFormatter -
	VcapServicesFormatter = `{
				"p-mongodb": [
				{"name": "%s","label": "p-mongodb",
				"tags": ["pivotal","mongodb"],"plan": "development",
          "credentials": {
            "uri": "%s",
            "scheme": "mongodb",
            "username": "c39642c7-xxxx-xxxx-xxxx-db67a3bbc98f",
            "password": "f6ac4b827xxxxxxxxxxx53533154e",
            "host": "192.168.88.184",
            "port": 27017,
            "database": "70ef645b-xxxx-4xxc-xxxx-94d5b0e5107f"
			}}]}`
	//VcapApplicationFormatter -
	VcapApplicationFormatter = `{
				"limits":{"mem":1024,"disk":1024,"fds":16384},
				"application_version":"56637561-e847-4023-87fa-1e476cb0b7e3",
				"application_name":"dispenserdev",
				"application_uris":["dispenserdev.cfapps.pez.pivotal.io","dispenserdev.pezapp.io"],
				"version":"56637561-e847-4023-87fa-1e476cb0b7e3",
				"name":"dispenserdev",
				"space_name":"pez-dev",
				"space_id":"ea88ed9e-91f1-4763-8eef-54fe38acf603",
				"uris":["dispenserdev.cfapps.pez.pivotal.io","dispenserdev.pezapp.io"],
				"users":null
			}`
)

var (
	server   *negroni.Negroni
	request  *http.Request
	recorder *httptest.ResponseRecorder
	appEnv   *cfenv.App
)

var _ = Describe("Server Integration Test", func() {

	BeforeEach(func() {
		os.Setenv("INVENTORY_DB_NAME", inventoryDB)
		os.Setenv("INVENTORY_DB_URI", inventoryURI)
		mongoURI := fmt.Sprintf(fakeURIFormatter, os.Getenv("MONGO_PORT_27017_TCP_ADDR"), os.Getenv("MONGO_PORT_27017_TCP_PORT"))
		os.Setenv("VCAP_SERVICES", fmt.Sprintf(VcapServicesFormatter, inventoryDB, mongoURI))
		os.Setenv("VCAP_APPLICATION", VcapApplicationFormatter)
		appEnv, _ = cfenv.Current()
	})

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

		XContext("when inventory exists", func() {
			It("returns a status code of 200 and a list of inventory", func() {
				server.ServeHTTP(recorder, request)
				Ω(recorder.Code).To(Equal(200))
				Ω(recorder.Body).To(ContainSubstring("2C.small"))
			})
		})
	})
})
