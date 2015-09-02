package pezinventory

import (
	"log"
	"net/http"

	"github.com/pivotal-pez/pezinventory/service/integrations"

	"gopkg.in/mgo.v2/bson"
)

//Inventory - inventory collection wrapper
type Inventory struct {
	ID         bson.ObjectId          `bson:"_id,omitempty"`
	SKU        string                 `json:"sku"`
	Tier       int                    `json:"tier"`
	Type       string                 `json:"type"`
	Size       string                 `json:"size"`
	Attributes map[string]interface{} `json:"attributes"`
	ItemStatus string                 `json:"item_status"`
	LeaseID    string                 `json:"lease_id"`
}

func listInventoryHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()
		log.Println("collection dialed successfully")

		//FIXME(dnem): short-circuited to insert dumb record
		// i := &Inventory{}
		// i.ID = bson.NewObjectId()
		// i.SKU = "2C.small"
		// i.Tier = 2
		// i.Type = "C"
		// i.Size = "small"

		// err = c.Insert(i)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		//result := Inventory{}
		// query := '{"sku": "2C.small"}'
		// err = c.Find(bson.M{"sku": "2C.small"}).One(&result)

		if cnt, err := collection.Count(); err == nil {
			log.Println("inventory count complete")
			Formatter().JSON(w, http.StatusOK, successMessage(cnt))
		} else {
			log.Println("inventory count failed")
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
		}
		// Formatter().JSON(w, http.StatusOK, successMessage(&result))
	}
}
