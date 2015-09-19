package pezinventory

import (
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/pivotal-pez/cfmgo"
	"github.com/unrolled/render"
)

var formatter *render.Render

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
	inventoryServiceURI := cfmgo.GetServiceBinding(inventoryServiceName, inventoryServiceURIName, appEnv)
	inventoryCollection := cfmgo.SetupDB(cfmgo.NewCollectionDialer, inventoryServiceURI, InventoryCollectionName)
	leaseCollection := cfmgo.SetupDB(cfmgo.NewCollectionDialer, inventoryServiceURI, LeaseCollectionName)

	n := negroni.Classic()
	mx := mux.NewRouter()

	mx.HandleFunc("/v1/inventory", ListInventoryItemsHandler(inventoryCollection)).Methods("GET")
	mx.HandleFunc("/v1/inventory", InsertInventoryItemHandler(inventoryCollection)).Methods("POST")
	mx.HandleFunc("/v1/leases/{id}", FindLeaseByIDHandler(leaseCollection)).Methods("GET")
	mx.HandleFunc("/v1/leases", FindLeasesHandler(leaseCollection)).Methods("GET")
	mx.HandleFunc("/v1/leases", LeaseInventoryItemHandler(inventoryCollection, leaseCollection)).Methods("POST")

	n.UseHandler(mx)

	return n
}
