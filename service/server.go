package pezinventory

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

//Render is a global reference to an initialized renderer
var Render *render.Render

func initRenderer() {
	Render = render.New(render.Options{
		IndentJSON: true,
	})
}

//NewServer configures and returns a Server.
//func NewServer(authHandler negroni.Handler) *negroni.Negroni {
func NewServer() *negroni.Negroni {

	Render = render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	//inventory routes
	mx.HandleFunc("/inventory", listInventory)
	n.UseHandler(mx)

	return n
}

//listInventory - controller function
func listInventory(w http.ResponseWriter, req *http.Request) {

	inventoryDB := os.Getenv("INVENTORY_DB_NAME")
	inventoryDBCollection := os.Getenv("INVENTORY_DB_COLLECTION")

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		Render.JSON(w, http.StatusInternalServerError, map[string]string{"error": "cannot connect to database"})
		return
	}
	defer session.Close()

	// query db
	c := session.DB(inventoryDB).C(inventoryDBCollection)

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

	result := Inventory{}
	err = c.Find(bson.M{"sku": "2C.small"}).One(&result)
	if err != nil {
		Render.JSON(w, http.StatusOK, map[string]string{"inventory": "[]"})
		log.Fatal(err)
	}

	// return results
	Render.JSON(w, http.StatusOK, &result)
}

//Inventory - inventory collection wrapper
type Inventory struct {
	ID         bson.ObjectId          `json:"_id"`
	SKU        string                 `json:"sku"`
	Tier       int                    `json:"tier"`
	Type       string                 `json:"type"`
	Size       string                 `json:"size"`
	Attributes map[string]interface{} `json:"attributes"`
	ItemStatus string                 `json:"item_status"`
	LeaseID    string                 `json:"lease_id"`
}
