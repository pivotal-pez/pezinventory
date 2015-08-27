package pezinventory

import (
	"net/http"

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
	Render.JSON(w, http.StatusOK, map[string]string{"inventory": "list"})
}

//Inventory - inventory collection wrapper
type Inventory struct {
	ID         string                 `json:"_id"`
	SKU        string                 `json:"sku"`
	Tier       int                    `json:"tier"`
	Type       string                 `json:"type"`
	Size       string                 `json:"size"`
	Attributes map[string]interface{} `json:"attributes"`
	ItemStatus string                 `json:"item_status"`
	LeaseID    string                 `json:"lease_id"`
}
