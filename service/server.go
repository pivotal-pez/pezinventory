package pezinventory

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

//Server wraps the Martini server struct
type Server *martini.ClassicMartini

//NewServer configures and returns a Server.
//TODO: Parameterize DB
func NewServer() Server {

	m := Server(martini.Classic())
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))

	//TODO: Database

	//types routes
	m.Group("/v1/types", func(r martini.Router) {
		ctrl := &TypeController{}
		r.Get("", ctrl.listTypes)
		r.Get("/:id", ctrl.getType)
		r.Get("/:id/items", ctrl.listTypeItems)
	})

	//items routes
	m.Group("/v1/items", func(r martini.Router) {
		ctrl := &ItemController{}
		r.Get("", ctrl.listItems)
		r.Get("/:id", ctrl.getItem)
		r.Get("/:id/history", ctrl.getItemHistory)
	})

	return m
}
