package pezinventory_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"gopkg.in/mgo.v2/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
	"github.com/pivotal-pez/pezinventory/service/fakes"
)

var _ = Describe("ListInventoryItemsHandler", func() {
	Context("when the handler is called", func() {
		var (
			fakeURI       = "mongodb://guid:quid@addr:port/guid"
			fakeInventory = []InventoryItem{
				InventoryItem{
					ID:                bson.NewObjectId(),
					SKU:               "2C.small",
					Tier:              2,
					OfferingType:      "C",
					Size:              "small",
					Status:            "available",
					PrivateAttributes: map[string]interface{}{"secret": "stuff"},
				},
				InventoryItem{
					ID:                bson.NewObjectId(),
					SKU:               "2C.small",
					Tier:              2,
					OfferingType:      "C",
					Size:              "small",
					Status:            "available",
					PrivateAttributes: map[string]interface{}{"secret": "stuff"},
				},
			}
			inventoryCollection = SetupDB(
				fakes.FakeNewCollectionDialer(fakeInventory),
				fakeURI,
				InventoryCollectionName)
		)

		It("should return true", func() {
			server := httptest.NewServer(http.HandlerFunc(http.HandlerFunc(ListInventoryItemsHandler(inventoryCollection))))
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

			Ω(payload).To(ContainSubstring("2C.small"))
			Ω(payload).ShouldNot(ContainElement("private_atrributes"))

		})
	})
})

//InsertInventoryItemHandler
var _ = Describe("InsertInventoryItemHandler", func() {
	var (
		fakeURI       = "mongodb://guid:quid@addr:port/guid"
		fakeInventory = []InventoryItem{}
		fakeItem      = &InventoryItem{
			ID:                bson.NewObjectId(),
			SKU:               "2C.small",
			Tier:              2,
			OfferingType:      "C",
			Size:              "small",
			Status:            "available",
			PrivateAttributes: map[string]interface{}{"secret": "stuff"},
		}
		inventoryCollection = SetupDB(
			fakes.FakeNewCollectionDialer(fakeInventory),
			fakeURI,
			InventoryCollectionName)
	)

	Context("when the handler is called", func() {
		server := httptest.NewServer(http.HandlerFunc(http.HandlerFunc(InsertInventoryItemHandler(inventoryCollection))))
		defer server.Close()

		b, err := json.Marshal(fakeItem)
		if err != nil {
			panic("cannot marshal InventoryItem")
		}

		res, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(b))
		if err != nil {
			log.Fatal(err)
		}

		payload, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		It("the return value should contain the item id", func() {
			Ω(payload).To(ContainSubstring(fakeItem.ID.Hex()))
		})
	})
})
