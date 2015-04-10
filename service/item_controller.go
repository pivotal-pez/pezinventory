package pezinventory

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

//ItemController - controller for searching item information
type ItemController struct {
}

//ListItems - returns a collection of Items
func (c *ItemController) listItems(render render.Render) {
	i := listItems()
	render.JSON(200, successMessage(i))
}

//GetItem - returns a single Item record
func (c *ItemController) getItem(params martini.Params, render render.Render) {
	i := getItem(params["id"])
	render.JSON(200, successMessage(i))
}

//GetItemHistory - returns the history for a given Item
func (c *ItemController) getItemHistory(params martini.Params, render render.Render) {
	render.JSON(501, errorMessage("Not Implemented"))
}
