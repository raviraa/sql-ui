package routes

import (
	"log"
	"sql-ui/controller"

	"github.com/gin-gonic/gin"
)

type (
	history struct {
		controller.Controller
	}

	historyData struct {
		Title string
	}
)

func (c *history) Get(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "main"
	page.Name = "history"

	log.Println(c.Container.Config.SaveConf())

	c.RenderPage(*ctx, page)
}
