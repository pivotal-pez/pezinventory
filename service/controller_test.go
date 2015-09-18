package pezinventory_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
)

var _ = Describe("ExtractQueryParams", func() {
	Context("when the handler is called with no query params", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			params := ExtractRequestParams(req.URL.Query())
			Formatter().JSON(w, http.StatusOK, &params)
			return
		}).Methods("GET")
		server := httptest.NewServer(mx)
		defer server.Close()

		url := server.URL + "/"
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		payload, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var rp RequestParams
		err = json.Unmarshal(payload, &rp)
		if err != nil {
			log.Fatal(err)
		}

		It("should have a default values", func() {
			Expect(rp.Selector).To(Equal(bson.M{}))
			Expect(rp.Scope).To(Equal(bson.M{}))
			Expect(rp.Limit).To(Equal(10))
			Expect(rp.Offset).To(Equal(0))
		})
	})
	Context("when the handler is called with parameters", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			// log.Printf("Scheme: %+v", req.URL.Scheme)
			// log.Printf("Host: %+v", req.URL.Host)
			// log.Printf("Path: %+v", req.URL.Path)
			// log.Printf("RawQuery: %+v", req.URL.RawQuery)
			// log.Printf("Values: %+v", req.URL.Query())

			params := ExtractRequestParams(req.URL.Query())
			Formatter().JSON(w, http.StatusOK, &params)
			return
		}).Methods("GET")
		server := httptest.NewServer(mx)
		defer server.Close()

		url := server.URL + "/?_id=1&limit=15&offset=30&scope=_id,status&status=available"
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		payload, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var rp RequestParams
		err = json.Unmarshal(payload, &rp)
		if err != nil {
			log.Fatal(err)
		}

		It("the request parameters object should be correctly populated", func() {
			Expect(rp.Selector["_id"]).NotTo(BeNil())
			Expect(rp.Selector["_id"].(string)).To(Equal("1"))
			Expect(rp.Selector["status"]).NotTo(BeNil())
			Expect(rp.Selector["status"].(string)).To(Equal("available"))
			Expect(rp.Scope["_id"]).NotTo(BeNil())
			Expect(rp.Scope["_id"].(float64)).To(Equal(float64(1)))
			Expect(rp.Scope["status"]).NotTo(BeNil())
			Expect(rp.Scope["status"].(float64)).To(Equal(float64(1)))
			Expect(rp.Limit).To(Equal(15))
			Expect(rp.Offset).To(Equal(30))
		})
	})

})
