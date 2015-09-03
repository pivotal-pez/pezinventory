package pezinventory

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pivotal-pez/pezinventory/service/integrations"

	"gopkg.in/mgo.v2/bson"
)

//InventoryItem - inventory collection wrapper
type InventoryItem struct {
	ID           bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	SKU          string                 `json:"sku"`
	Tier         int                    `json:"tier"`
	OfferingType string                 `json:"offeringType"`
	Size         string                 `json:"size"`
	Attributes   map[string]interface{} `json:"attributes"`
	Status       string                 `json:"status"`
	LeaseID      string                 `json:"lease_id"`
}

func listInventoryItemsHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		//FIXME(dnem) short-circuited to return all; update to ingest request params
		selector := bson.M{}
		items := make([]InventoryItem, 0)

		if err := collection.Find(selector, &items); err == nil {
			Formatter().JSON(w, http.StatusOK, successMessage(&items))
		} else {
			log.Println("inventory find failed")
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
		}
	}
}

func insertInventoryItemHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var i InventoryItem
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&i)
		if err != nil {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage(err.Error()))
			return
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
