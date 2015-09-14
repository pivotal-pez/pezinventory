package pezinventory_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
	"github.com/pivotal-pez/pezinventory/service/fakes"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("FindLeaseByIDHandler", func() {
	var (
		id1                 = bson.NewObjectId()
		id2                 = bson.NewObjectId()
		fakeURI             = "mongodb://guid:guid@addr:port/guid"
		fakeLeaseCollection = []Lease{
			Lease{
				ID:                id1,
				InventoryItemID:   bson.NewObjectId(),
				User:              "testuser-1",
				Duration:          "14 days",
				StartDate:         time.Now().AddDate(0, 0, -2).String(),
				EndDate:           time.Now().AddDate(0, 0, 11).String(),
				Status:            "active",
				Attributes:        map[string]interface{}{"public": "info"},
				PrivateAttributes: map[string]interface{}{"secret": "stuff"},
			},
			Lease{
				ID:                id2,
				InventoryItemID:   bson.NewObjectId(),
				User:              "testuser-2",
				Duration:          "28 days",
				StartDate:         time.Now().AddDate(0, 0, -14).String(),
				EndDate:           time.Now().AddDate(0, 0, 13).String(),
				Status:            "active",
				Attributes:        map[string]interface{}{"public": "info"},
				PrivateAttributes: map[string]interface{}{"secret": "stuff"},
			},
		}
		leaseCollection = SetupDB(
			fakes.FakeNewCollectionDialer(fakeLeaseCollection),
			fakeURI,
			LeaseCollectionName)
	)

	Context("when the hander is called without a LeaseID", func() {
		It("should return an error response", func() {
			server := httptest.NewServer(http.HandlerFunc(http.HandlerFunc(FindLeaseByIDHandler(leaseCollection))))
			defer server.Close()

			res, err := http.Get(server.URL)
			if err != nil {
				log.Fatal(err)
			}

			payload, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			Ω(payload).To(ContainSubstring("error"))
			Ω(payload).Should(ContainSubstring("lease id must be specified"))
			Ω(payload).ShouldNot(ContainElement("data"))
		})
	})

	Context("when the hander is called with a valid LeaseID for record 1", func() {
		It("should return lease record 1", func() {

			mx := mux.NewRouter()
			mx.HandleFunc("/v1/leases/{id}", FindLeaseByIDHandler(leaseCollection)).Methods("GET")

			server := httptest.NewServer(mx)
			defer server.Close()

			leaseURL := server.URL + "/v1/leases/0"

			res, err := http.Get(leaseURL)
			if err != nil {
				log.Fatal(err)
			}

			payload, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			Ω(payload).To(ContainSubstring("success"))
			Ω(payload).To(ContainSubstring(id1.Hex()))
		})
	})

	Context("when the hander is called with a valid LeaseID for record 2", func() {
		It("should return lease record 2", func() {

			mx := mux.NewRouter()
			mx.HandleFunc("/v1/leases/{id}", FindLeaseByIDHandler(leaseCollection)).Methods("GET")

			server := httptest.NewServer(mx)
			defer server.Close()

			leaseURL := server.URL + "/v1/leases/1"

			res, err := http.Get(leaseURL)
			if err != nil {
				log.Fatal(err)
			}

			payload, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			Ω(payload).To(ContainSubstring("success"))
			Ω(payload).To(ContainSubstring(id2.Hex()))
		})
	})

})
