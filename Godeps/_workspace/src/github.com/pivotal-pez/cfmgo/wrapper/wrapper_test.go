package wrapper_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/cfmgo/params"
	. "github.com/pivotal-pez/cfmgo/wrapper"
)

var _ = Describe("ResponseWrapper", func() {
	Context("when the message is wrapped in ErrorWrapper", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			err := errors.New("must feed a hamburger to the gnome before continuing")
			Formatter().JSON(w, http.StatusOK, ErrorWrapper(err.Error()))
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

		var rw ResponseWrapper
		err = json.Unmarshal(payload, &rw)
		if err != nil {
			log.Fatal(err)
		}

		It("should have a default values", func() {
			Expect(rw.Status).To(Equal("error"))
			Expect(rw.Message).To(Equal("must feed a hamburger to the gnome before continuing"))
		})
	})

	Context("when the is wrapped with SuccessWrapper", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			params := ExtractRequestParams(req.URL.Query())
			Formatter().JSON(w, http.StatusOK, SuccessWrapper(&params))
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

		var rw ResponseWrapper
		err = json.Unmarshal(payload, &rw)
		if err != nil {
			log.Fatal(err)
		}

		It("should have a default values", func() {
			Expect(rw.Status).To(Equal("success"))
			Expect(rw.Data).NotTo(BeNil())
			Expect(rw.Count).To(Equal(0))
			Expect(rw.Message).To(Equal(""))
		})
	})

	Context("when the is wrapped with CollectionWrapper", func() {
		mx := mux.NewRouter()
		mx.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			params := ExtractRequestParams(req.URL.Query())
			count := 67
			Formatter().JSON(w, http.StatusOK, CollectionWrapper(params, count))
			return
		}).Methods("GET")
		server := httptest.NewServer(mx)
		defer server.Close()

		url := server.URL + "/?limit=2&offset=1&scope=fluffy&a=1&b=2"
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		payload, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var rw ResponseWrapper
		err = json.Unmarshal(payload, &rw)
		if err != nil {
			log.Fatal(err)
		}

		b, err := json.Marshal(rw.Data)
		if err != nil {
			log.Fatal(err)
		}

		var p RequestParams
		err = json.Unmarshal(b, &p)
		if err != nil {
			log.Fatal(err)
		}

		It("should have a default values", func() {
			Expect(rw.Status).To(Equal("success"))
			Expect(rw.Data).NotTo(BeNil())
			Expect(rw.Count).To(Equal(67))
			Expect(rw.Message).To(Equal(""))
			Expect(p.Limit()).To(Equal(2))
			Expect(p.Offset()).To(Equal(1))
			Expect(p.Scope()["fluffy"]).To(Equal(float64(1)))
			Expect(p.Selector()["a"]).To(Equal("1"))
			Expect(p.Selector()["b"]).To(Equal("2"))
		})
	})
})
