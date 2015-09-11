package pezinventory

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/pivotal-pez/pezinventory/service/integrations"

	"github.com/unrolled/render"
)

var formatter *render.Render = nil

//Formatter returns the address for a global response formatter
//realized in the `github.com/unrolled/render` package.
func Formatter() *render.Render {
	if formatter == nil {
		formatter = render.New(render.Options{
			IndentJSON: true,
		})
	}
	return formatter
}

//NewServer configures and returns a Server.
func NewServer(appEnv *cfenv.App) *negroni.Negroni {

	inventoryServiceName := os.Getenv("INVENTORY_DB_NAME")
	inventoryServiceURIName := os.Getenv("INVENTORY_DB_URI")
	inventoryServiceURI := getServiceBinding(inventoryServiceName, inventoryServiceURIName, appEnv)
	inventoryCollection := SetupDB(integrations.NewCollectionDialer, inventoryServiceURI, InventoryCollectionName)
	leaseServiceName := os.Getenv("LEASE_DB_NAME")
	leaseServiceURIName := os.Getenv("LEASE_DB_URI")
	leaseServiceURI := getServiceBinding(leaseServiceName, leaseServiceURIName, appEnv)
	leaseCollection := SetupDB(integrations.NewCollectionDialer, leaseServiceURI, LeaseCollectionName)

	n := negroni.Classic()
	mx := mux.NewRouter()

	mx.HandleFunc("/v1/inventory", ListInventoryItemsHandler(inventoryCollection)).Methods("GET")
	mx.HandleFunc("/v1/inventory", InsertInventoryItemHandler(inventoryCollection)).Methods("POST")
	mx.HandleFunc("/v1/lease/{id}", FindLeaseByIDHandler(leaseCollection)).Methods("GET")
	mx.HandleFunc("/v1/lease", InsertLeaseRecordHandler(leaseCollection)).Methods("POST")

	n.UseHandler(mx)

	return n
}

func getServiceBinding(serviceName string, serviceURIName string, appEnv *cfenv.App) (serviceURI string) {

	if service, err := appEnv.Services.WithName(serviceName); err == nil {
		if serviceURI = service.Credentials[serviceURIName].(string); serviceURI == "" {
			panic(fmt.Sprintf("we pulled an empty connection string %s from %v - %v", serviceURI, service, service.Credentials))
		}

	} else {
		panic(fmt.Sprint("Experienced an error trying to grab service binding information:", err.Error()))
	}
	return
}

func SetupDB(dialer integrations.CollectionDialer, URI string, collectionName string) (collection integrations.Collection) {
	var (
		err      error
		dialInfo *mgo.DialInfo
	)

	if dialInfo, err = mgo.ParseURL(URI); err != nil || dialInfo.Database == "" {
		panic(fmt.Sprintf("can not parse given URI %s due to error: %s", URI, err.Error()))
	}

	if collection, err = dialer(URI, dialInfo.Database, collectionName); err != nil {
		panic(fmt.Sprintf("can not dial connection due to error: %s URI:%s col:%s db:%s", err.Error(), URI, collectionName, dialInfo.Database))
	}
	return
}
