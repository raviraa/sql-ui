package routes

import (
	"sql-ui/controller"

	"github.com/gin-gonic/gin"
)

type (
	home struct {
		controller.Controller
	}

	post struct {
		Title string
		Body  string
	}
)

func (c *home) Get(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "main"
	page.Name = "home"

	c.RenderPage(*ctx, page)
}
