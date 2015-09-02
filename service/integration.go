package pezinventory

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
)

//IntegrationWrapper holds a set of named Integration objects
type IntegrationWrapper map[string]*Integration

//Integration contains integrated service details.
type Integration struct {
	URI        string
	DB         string
	Collection string
}

func NewIntegrationWrapper() IntegrationWrapper {
	return make(map[string]*Integration)
}

func GetIntegrations() (iw IntegrationWrapper) {
	var inventoryServiceURI string
	inventoryDB := os.Getenv("INVENTORY_DB_NAME")
	inventoryCredsURIName := os.Getenv("INVENTORY_DB_URI")
	inventoryDBCollection := os.Getenv("INVENTORY_DB_COLLECTION")

	appEnv, _ := cfenv.Current()
	if inventoryService, err := appEnv.Services.WithName(inventoryDB); err == nil {
		if inventoryServiceURI = inventoryService.Credentials[inventoryCredsURIName].(string); inventoryServiceURI == "" {
			panic(fmt.Sprintf("Retrieved empty connection string %s from %v - %v", inventoryServiceURI, inventoryService, inventoryService.Credentials))
		}
	} else {
		panic(fmt.Sprint("Unable to retrieve service binding information", err.Error()))
	}

	inventory := Integration{
		inventoryServiceURI,
		inventoryDB,
		inventoryDBCollection,
	}

	iw = NewIntegrationWrapper()
	iw["inventory"] = &inventory

	return
}
