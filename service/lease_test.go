package pezinventory_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/dnem/paged"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-pez/cfmgo"
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
				DurationDays:      14,
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
				DurationDays:      28,
				StartDate:         time.Now().AddDate(0, 0, -14).String(),
				EndDate:           time.Now().AddDate(0, 0, 13).String(),
				Status:            "active",
				Attributes:        map[string]interface{}{"public": "info"},
				PrivateAttributes: map[string]interface{}{"secret": "stuff"},
			},
		}
		leaseCollection = cfmgo.Connect(
			fakes.FakeNewCollectionDialer(fakeLeaseCollection),
			fakeURI,
			LeaseCollectionName)
	)

	Context("when the FindLeasesHandler is called without parameters", func() {
		It("should return a collection of 2 lease objects", func() {
			mx := mux.NewRouter()
			mx.HandleFunc("/v1/leases", FindLeasesHandler(leaseCollection)).Methods("GET")

			server := httptest.NewServer(mx)
			defer server.Close()

			leaseURL := server.URL + "/v1/leases"

			res, err := http.Get(leaseURL)
			if err != nil {
				log.Fatal(err)
			}

			payload, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			var rw paged.ResponseWrapper
			err = json.Unmarshal(payload, &rw)
			if err != nil {
				log.Fatal(err)
			}

			Expect(rw.Status).To(Equal("success"))
			Expect(rw.Count).To(Equal(2))
			Expect(payload).To(ContainSubstring("testuser-1"))
			Expect(payload).To(ContainSubstring("testuser-2"))
		})
	})

	Context("when the handler is called without a LeaseID", func() {
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

			Expect(payload).To(ContainSubstring("error"))
			Expect(payload).Should(ContainSubstring("lease id must be specified"))
			Expect(payload).ShouldNot(ContainElement("data"))
		})
	})

	Context("when the handler is called with a valid LeaseID for record 1", func() {
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

			Expect(payload).To(ContainSubstring("success"))
			Expect(payload).To(ContainSubstring(id1.Hex()))
		})
	})

	Context("when the handler is called with a valid LeaseID for record 2", func() {
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

			Expect(payload).To(ContainSubstring("success"))
			Expect(payload).To(ContainSubstring(id2.Hex()))
		})
	})

})
