package pezinventory

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

//TypeController - controller for searching type information
type TypeController struct {
}

func (c *TypeController) listTypes(render render.Render) {
	t := listTypes()
	render.JSON(200, successMessage(t))
}

func (c *TypeController) getType(params martini.Params, render render.Render) {
	t := getType(params["id"])
	render.JSON(200, successMessage(t))
}

func (c *TypeController) listTypeItems(params martini.Params, render render.Render) {
	i := listItemsByType(params["id"])
	render.JSON(200, successMessage(i))
}
