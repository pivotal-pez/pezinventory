package pezinventory

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/cfmgo/params"
	"github.com/pivotal-pez/cfmgo/wrapper"
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

// @Title ListInventoryItems
// @Description Returns a collection of inventory items based on the supplied parameters
// @Accept json
// @Param  limit    query   string   false   "Maximum number of records in result set"
// @Param  offset   query   string   false   "Number of records to skip; supports paging"
// @Param  scope    query   string   false   "Comma-separated list of fields to return"
// @Success 200 {object} wrapper.ResponseWrapper
// @Failure 404 {object} wrapper.ResponseWrapper
// @Resource /v1/inventory
// @Router /v1/inventory [get]
//ListInventoryItemsHandler returns a collection of InventoryItems based on supplied paramaters.
func ListInventoryItemsHandler(collection cfmgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		params := params.Extract(req.URL.Query())

		items := make([]RedactedInventoryItem, 0)

		if count, err := collection.Find(params, &items); err == nil {
			Formatter().JSON(w, http.StatusOK, wrapper.Collection(&items, count))
		} else {
			Formatter().JSON(w, http.StatusNotFound, wrapper.Error(err.Error()))
		}
	}
}

// @Title Insert / Update Inventory Item
// @Description Uses MongoDB UPSERT to add new records, or update existing records, to the InventoryItems collection.
// @Accept json
// @Param  body body string true "JSON string of InventoryItem Record"
// @Success 200 {object} wrapper.ResponseWrapper
// @Failure 500 {object} wrapper.ResponseWrapper
// @Resource /v1/inventory
// @Router /v1/inventory [post]
//InsertInventoryItemHandler uses MongoDB UPSERT to add new records, or update existing records,
//to the InventoryItems collection.
func InsertInventoryItemHandler(collection cfmgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var i InventoryItem
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&i)
		if err != nil {
			Formatter().JSON(w, http.StatusBadRequest, wrapper.Error(err.Error()))
			return
		}

		if i.ID == "" {
			i.ID = bson.NewObjectId()
		}

		collection.Wake()
		info, err := collection.UpsertID(i.ID, i)
		if err != nil {
			log.Println("could not create InventoryItem record")
			Formatter().JSON(w, http.StatusInternalServerError, wrapper.Error(err.Error()))
		} else {
			Formatter().JSON(w, http.StatusOK, wrapper.One(info))
		}
	}
}

//InventoryItemReservingStatus updates the status from "available" to "reserving".
func InventoryItemReservingStatus(id bson.ObjectId, collection cfmgo.Collection) error {
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
func InventoryItemAvailableStatus(id bson.ObjectId, collection cfmgo.Collection) error {
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
func InventoryItemLeasedStatus(id bson.ObjectId, leaseId bson.ObjectId, collection cfmgo.Collection) error {
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
