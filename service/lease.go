package pezinventory

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pivotal-pez/pezinventory/service/integrations"

	"gopkg.in/mgo.v2/bson"
)

//Lease wraps the leases collection.
type Lease struct {
	ID                bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	InventoryItemID   bson.ObjectId          `bson:"inventory_item_id,omitempty" json:"inventory_item_id"`
	User              string                 `json:"user"`
	Duration          string                 `json:"duration"`
	StartDate         string                 `json:"start_date"`
	EndDate           string                 `json:"end_date"`
	Status            string                 `json:"status"`
	Attributes        map[string]interface{} `json:"attributes"`
	PrivateAttributes map[string]interface{} `json:"private_attributes,omitempty"`
}

//RedactedLease wraps the leases collection omitting private attributes.
type RedactedLease struct {
	ID              bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	InventoryItemID bson.ObjectId          `bson:"inventory_item_id,omitempty" json:"inventory_item_id"`
	User            string                 `json:"user"`
	Duration        string                 `json:"duration"`
	StartDate       string                 `json:"start_date"`
	EndDate         string                 `json:"end_date"`
	Status          string                 `json:"status"`
	Attributes      map[string]interface{} `json:"attributes"`
}

//FindLeaseByIDHandler will return a redacted lease record for the given ID.
func FindLeaseByIDHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		id := mux.Vars(req)["id"]
		if id == "" {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage("lease id must be specified"))
			return
		}

		lease := new(RedactedLease)
		if err := collection.FindOne(id, lease); err == nil {
			Formatter().JSON(w, http.StatusOK, successMessage(lease))
		} else {
			log.Println("lease lookup failed")
			Formatter().JSON(w, http.StatusNotFound, errorMessage(err.Error()))
		}
	}
}

//LeaseInventoryItemHandler creates a new lease record against an available InventoryItem
//and calls dispenser to provision that InventoryItem to the requestor.
//NOTE: call to dispenser not implemented
func LeaseInventoryItemHandler(ic integrations.Collection, lc integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var leaseObj Lease
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&leaseObj)
		if err != nil {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage(err.Error()))
			return
		}

		if leaseObj.InventoryItemID == "" {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage("inventory_item_id must be specified"))
			return
		}

		err = InventoryItemReservingStatus(leaseObj.InventoryItemID, ic)
		if err != nil {
			Formatter().JSON(w, http.StatusNotFound, errorMessage(err.Error()))
			return
		}

		lc.Wake()
		leaseObj.ID = bson.NewObjectId()
		_, err = lc.UpsertID(leaseObj.ID, leaseObj)
		if err != nil {
			e := InventoryItemAvailableStatus(leaseObj.InventoryItemID, ic)
			if e != nil {
				log.Printf("Could not release Inventory Item %s", leaseObj.InventoryItemID.Hex())
			}
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
			return
		}

		err = InventoryItemLeasedStatus(leaseObj.InventoryItemID, leaseObj.ID, ic)
		if err != nil {
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
			return
		}

		Formatter().JSON(w, http.StatusOK, leaseObj)
	}
}

//InsertLeaseRecordHandler performs an upsert on a new/existing lease record.
//FIXME(dnem) This should be modified to handle lease updates via  PATCH /v1/leases
func InsertLeaseRecordHandler(collection integrations.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var obj Lease
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&obj)
		if err != nil {
			Formatter().JSON(w, http.StatusBadRequest, errorMessage(err.Error()))
			return
		}

		if obj.ID == "" {
			obj.ID = bson.NewObjectId()
		}

		collection.Wake()
		info, err := collection.UpsertID(obj.ID, obj)
		if err != nil {
			log.Println("could not create Lease record")
			Formatter().JSON(w, http.StatusInternalServerError, errorMessage(err.Error()))
		} else {
			log.Println(info)
			Formatter().JSON(w, http.StatusOK, successMessage(info))
		}
	}
}
