package routes

import (
	"github.com/raviraa/sql-ui/config"
	"github.com/raviraa/sql-ui/controller"

	"github.com/gin-gonic/gin"
)

type (
	history struct {
		controller.Controller
	}

	historyData struct {
		Entries []config.HistEntry
	}
)

func (c *history) Get(ctx *gin.Context) {
	page := controller.NewPage(ctx, c.Container)
	page.Layout = "main"
	page.Name = "history"
	page.Data = historyData{
		Entries: c.Container.Config.GetHistEntries(),
	}

	c.RenderPage(ctx, page)
}
