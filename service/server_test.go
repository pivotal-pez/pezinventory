package pezinventory_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/pezinventory/service"
	"github.com/pivotal-pez/pezinventory/service/fakes"
)

var _ = Describe("listInventoryItemsController", func() {
	Context("when the handler is called", func() {
		var (
			fakeURI       = "mongodb://guid:quid@addr:port/guid"
			fakeInventory = []InventoryItem{
				InventoryItem{
					ID:           "abcdef1",
					SKU:          "2C.small",
					Tier:         2,
					OfferingType: "C",
					Size:         "small",
					Status:       "available",
				},
				InventoryItem{
					ID:           "abcdef2",
					SKU:          "2C.small",
					Tier:         2,
					OfferingType: "C",
					Size:         "small",
					Status:       "available",
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
		})
	})
})
