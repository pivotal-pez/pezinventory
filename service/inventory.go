package pezinventory

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dnem/paged"
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
	ID           bson.ObjectId          `bson:"_id,omitempty" json:"id,omitempty"`
	SKU          string                 `json:"sku,omitempty"`
	Tier         int                    `json:"tier,omitempty"`
	OfferingType string                 `json:"offering_type,omitempty"`
	Size         string                 `json:"size,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	Status       string                 `json:"status,omitempty"`
	LeaseID      bson.ObjectId          `bson:"lease_id,omitempty" json:"lease_id,omitempty"`
}

//ListInventoryItemsHandler -
func ListInventoryItemsHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		params := paged.ExtractRequestParams(req.URL.Query())

		items := make([]RedactedInventoryItem, 0)

		if count, err := collection.Find(params.Selector(), params.Scope(), params.Limit(), params.Offset(), &items); err == nil {
			Formatter().JSON(w, http.StatusOK, paged.CollectionWrapper(&items, count))
		} else {
			Formatter().JSON(w, http.StatusNotFound, paged.ErrorWrapper(err.Error()))
		}
	}
}

//InsertInventoryItemHandler -
//FIXME(dnem) consider returning ID rather than mgo.ChangeInfo
func InsertInventoryItemHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var i InventoryItem
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&i)
		if err != nil {
			Formatter().JSON(w, http.StatusBadRequest, paged.ErrorWrapper(err.Error()))
			return
		}

		if i.ID == "" {
			i.ID = bson.NewObjectId()
		}

		collection.Wake()
		info, err := collection.UpsertID(i.ID, i)
		if err != nil {
			log.Println("could not create InventoryItem record")
			Formatter().JSON(w, http.StatusInternalServerError, paged.ErrorWrapper(err.Error()))
		} else {
			Formatter().JSON(w, http.StatusOK, paged.SuccessWrapper(info))
		}
	}
}

//InventoryItemReservingStatus updates the status from "available" to "reserving".
func InventoryItemReservingStatus(id bson.ObjectId, collection integrations.Collection) error {
	var obj RedactedInventoryItem

	sel := bson.M{
		"_id":    id,
		"status": InventoryItemStatusAvailable,
	}

	update := bson.M{
		"$set": bson.M{
			"status": InventoryItemStatusReserving,
		},
	}

	collection.Wake()
	_, err := collection.FindAndModify(sel, update, &obj)
	if err != nil {
		err = ErrInventoryNotAvailable
	}
	return err
}

//InventoryItemAvailableStatus reverts the status from "reserving" to "available" in the case
//where a lease operation is unsuccessful.
func InventoryItemAvailableStatus(id bson.ObjectId, collection integrations.Collection) error {
	var obj RedactedInventoryItem

	sel := bson.M{
		"_id":    id,
		"status": InventoryItemStatusReserving,
	}

	update := bson.M{
		"$set": bson.M{
			"status": InventoryItemStatusAvailable,
		},
	}

	collection.Wake()
	_, err := collection.FindAndModify(sel, update, &obj)
	return err
}

//InventoryItemLeasedStatus updates the status from "reserving" to "leased" and supplies
//the lease_id value.
func InventoryItemLeasedStatus(id bson.ObjectId, leaseId bson.ObjectId, collection integrations.Collection) error {
	var obj RedactedInventoryItem

	sel := bson.M{
		"_id":    id,
		"status": InventoryItemStatusReserving,
	}

	update := bson.M{
		"$set": bson.M{
			"status":   InventoryItemStatusLeased,
			"lease_id": leaseId,
		},
	}

	collection.Wake()
	_, err := collection.FindAndModify(sel, update, &obj)
	return err
}
