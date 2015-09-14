package pezinventory

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pivotal-pez/pezinventory/service/integrations"

	"gopkg.in/mgo.v2/bson"
)

//InventoryItem wraps the inventory collection.
type InventoryItem struct {
	ID                bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	SKU               string                 `json:"sku"`
	Tier              int                    `json:"tier"`
	OfferingType      string                 `json:"offering_type"`
	Size              string                 `json:"size"`
	Attributes        map[string]interface{} `json:"attributes"`
	PrivateAttributes map[string]interface{} `json:"private_attributes,omitempty"`
	Status            string                 `json:"status"`
	LeaseID           bson.ObjectId          `bson:"lease_id,omitempty" json:"lease_id"`
}

//RedactedInventoryItem wraps the inventory collection omitting private attributes.
type RedactedInventoryItem struct {
	ID           bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	SKU          string                 `json:"sku"`
	Tier         int                    `json:"tier"`
	OfferingType string                 `json:"offering_type"`
	Size         string                 `json:"size"`
	Attributes   map[string]interface{} `json:"attributes"`
	Status       string                 `json:"status"`
	LeaseID      bson.ObjectId          `bson:"lease_id,omitempty" json:"lease_id"`
}

//ListInventoryItemsHandler -
// currently selects using a nil selector
func ListInventoryItemsHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		selector := bson.M{}
		items := make([]RedactedInventoryItem, 0)

		if err := collection.Find(selector, &items); err == nil {
			Formatter().JSON(w, http.StatusOK, successMessage(&items))
		} else {
			log.Println("inventory find failed")
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
		}
	}
}

//InsertInventoryItemHandler -
func InsertInventoryItemHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var i InventoryItem
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&i)
		if err != nil {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage(err.Error()))
		} else {
			i.ID = bson.NewObjectId()
		}

		collection.Wake()
		info, err := collection.UpsertID(i.ID, i)
		if err != nil {
			log.Println("could not create InventoryItem record")
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
		} else {
			log.Println(info)
			//FIXME(dnem) consider returning ID rather than mgo.ChangeInfo
			Formatter().JSON(w, http.StatusOK, successMessage(info))
		}
	}
}
