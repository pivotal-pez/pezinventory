package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"gopkg.in/mgo.v2/bson"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-pez/cfmgo/wrap"
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
				"application_name":"inventorydev",
				"application_uris":["inventorydev.cfapps.core.pez.pivotal.io"],
				"version":"56637561-e847-4023-87fa-1e476cb0b7e3",
				"name":"inventorydev",
				"space_name":"pez-dev",
				"space_id":"ea88ed9e-91f1-4763-8eef-54fe38acf603",
				"uris":["inventorydev.cfapps.core.pez.pivotal.io"],
				"users":null
			}`
)

var (
	server     *negroni.Negroni
	request    *http.Request
	recorder   *httptest.ResponseRecorder
	appEnv     *cfenv.App
	mongoURI   string
	newLeaseID bson.ObjectId
	item       = &InventoryItem{
		ID:                bson.NewObjectId(),
		SKU:               "2C.small",
		Tier:              2,
		OfferingType:      "C",
		Size:              "small",
		Status:            "available",
		PrivateAttributes: map[string]interface{}{"secret": "stuff"},
	}
	lease = &Lease{
		InventoryItemID: item.ID,
		User:            "testuser",
		Status:          "active",
	}
)

var _ = Describe("Server Integration Tests", func() {

	BeforeEach(func() {
		os.Setenv("INVENTORY_DB_NAME", inventoryDB)
		os.Setenv("INVENTORY_DB_URI", inventoryURI)
		mongoURI = fmt.Sprintf(fakeURIFormatter,
			os.Getenv("MONGO_PORT_27017_TCP_ADDR"),
			os.Getenv("MONGO_PORT_27017_TCP_PORT"))
		os.Setenv("VCAP_SERVICES",
			fmt.Sprintf(VcapServicesFormatter,
				inventoryDB, mongoURI))
		os.Setenv("VCAP_APPLICATION", VcapApplicationFormatter)
		appEnv, _ = cfenv.Current()
	})

	Context("when inventory does not exist", func() {
		BeforeEach(func() {
			server = NewServer(appEnv)
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest("GET", "/v1/inventory", nil)
		})

		It("GET /v1/inventory returns a status code of 200 and an empty set", func() {
			server.ServeHTTP(recorder, request)
			Ω(recorder.Code).To(Equal(200))
			Ω(recorder.Body).To(ContainSubstring("[]"))
		})
	})

	Context("when inventory is added", func() {
		BeforeEach(func() {
			b, err := json.Marshal(item)
			if err != nil {
				log.Fatal(err)
			}

			server = NewServer(appEnv)
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest(
				"POST",
				"/v1/inventory",
				bytes.NewBuffer(b))
		})

		It("POST /v1/inventory returns a status code of 200 and contains the id", func() {
			server.ServeHTTP(recorder, request)
			Ω(recorder.Code).To(Equal(200))
			Ω(recorder.Body).To(ContainSubstring(item.ID.Hex()))
		})

		It("GET /v1/inventory returns a status code of 200 and a list of inventory", func() {
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest("GET", "/v1/inventory", nil)
			server.ServeHTTP(recorder, request)
			Ω(recorder.Code).To(Equal(200))
			Ω(recorder.Body).To(ContainSubstring("2C.small"))
		})

	})

	Context("when inventory is available to lease", func() {
		BeforeEach(func() {
			b, err := json.Marshal(lease)
			if err != nil {
				log.Fatal(err)
			}

			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest(
				"POST",
				"/v1/leases",
				bytes.NewBuffer(b))
		})

		It("POST /v1/leases returns a status code of 200 and the new lease record", func() {
			var rm wrap.ResponseWrapper
			var newLease RedactedLease
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))
			err := json.Unmarshal(recorder.Body.Bytes(), &rm)
			if err != nil {
				log.Fatal(err)
			}
			Expect(rm.Status).To(ContainSubstring("success"))
			b, err := json.Marshal(rm.Data)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(b, &newLease)
			if err != nil {
				log.Fatal(err)
			}
			newLeaseID = newLease.ID
			Expect(newLease.User).To(ContainSubstring("testuser"))
		})
	})

	Context("when a lease exists", func() {
		BeforeEach(func() {
			server = NewServer(appEnv)
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest(
				"GET",
				"/v1/leases/"+newLeaseID.Hex(),
				nil)
		})

		It("GET /v1/leases/:id returns a lease object", func() {
			server.ServeHTTP(recorder, request)
			Ω(recorder.Code).To(Equal(200))
			Ω(recorder.Body).To(ContainSubstring("testuser"))
		})
	})

	Context("when leases exist", func() {
		BeforeEach(func() {
			server = NewServer(appEnv)
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest(
				"GET",
				"/v1/leases?user=testuser",
				nil)
		})

		It("GET /v1/leases?user=testuser returns a collection of lease objects", func() {
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))
			Expect(recorder.Body).To(ContainSubstring("testuser"))
			Expect(recorder.Body).To(ContainSubstring("count"))
		})
	})

	Context("when inventory is not available to lease", func() {
		BeforeEach(func() {
			b, err := json.Marshal(lease)
			if err != nil {
				log.Fatal(err)
			}

			server = NewServer(appEnv)
			recorder = httptest.NewRecorder()
			request, _ = http.NewRequest(
				"POST",
				"/v1/leases",
				bytes.NewBuffer(b))
		})

		It("POST /v1/leases returns a status code of 200 and the new lease record", func() {
			server.ServeHTTP(recorder, request)
			Ω(recorder.Code).To(Equal(404))
			Ω(recorder.Body).To(ContainSubstring("The InventoryItem specified is not available"))
		})
	})
})
