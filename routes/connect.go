package routes

import (
	"log"
	"sql-ui/controller"
	"sql-ui/services/qrunner"

	"github.com/gin-gonic/gin"
)

type (
	connect struct {
		controller.Controller
	}
	connectData struct {
		Errmsg       string
		ConnectedDSN string
		DSN          string
		DSNsRecent   []string
	}
)

func (c *connect) Get(ctx *gin.Context) {
	page := controller.NewPage(ctx, c.Container)
	data := connectData{}
	page.Layout = "main"
	page.Name = "connect"
	if c.Container.Qrunner != nil {
		data.ConnectedDSN = "connected" + c.Container.Qrunner.Dsn()
	}
	data.DSNsRecent = c.Container.Config.DSNsRecent()
	page.Data = data
	c.RenderPage(ctx, page)
}

func (c *connect) Post(ctx *gin.Context) {
	data := connectData{}
	page := controller.NewPage(ctx, c.Container)
	page.Name = "connect"
	dsn := ctx.PostForm("dsn")
	data.DSN = dsn
	data.DSNsRecent = c.Container.Config.DSNsRecent()

	if ctx.PostForm("disconnect") == "true" {
		data.Errmsg = "Disconnected"
		if c.Container.Qrunner != nil {
			c.Container.Qrunner.Close()
			c.Container.Qrunner = nil
		}
	} else if dsn != "" {
		qr, err := qrunner.New(dsn)
		if err != nil {
			log.Println("Connect failed ", err)
			data.Errmsg = err.Error()
		} else {
			c.Container.Qrunner = qr
			c.Container.Config.AddDSN(dsn)
			data.ConnectedDSN = "connected. " + qr.Dsn()
		}
	}
	page.Data = data
	c.RenderPage(ctx, page)
}
