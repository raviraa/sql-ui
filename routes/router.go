package routes

import (
	"github.com/raviraa/sql-ui/config"
	"github.com/raviraa/sql-ui/controller"
	"github.com/raviraa/sql-ui/middleware"
	"github.com/raviraa/sql-ui/services"
	"github.com/raviraa/sql-ui/services/qrunner"
	"github.com/raviraa/sql-ui/static"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuildRouter(c *services.Container) {
	c.Web.Group("").
		Use(middleware.CacheControl(config.CacheControl)).
		StaticFS(config.StaticPrefix, http.FS(static.StaticFS))

	g := c.Web.Group("")
	g.Use(
		gin.ErrorLogger(),
		gin.Recovery(),
	)

	// Base controller
	ctr := controller.NewController(c)
	userRoutes(c, g, ctr)

}

func userRoutes(c *services.Container, g *gin.RouterGroup, ctr controller.Controller) {

	query := query{Controller: ctr}
	g.GET("/", query.Get)
	g.POST("/query", query.Post)

	tables := tables{ctr}
	g.GET("/tables/meta/:metacmd", tables.Meta)
	g.GET("/tables/browse", tables.Browse)

	history := history{ctr}
	g.GET("/history", history.Get)

	if config.DbDsn != "" {
		// forced db connect
		var err error
		c.Qrunner, err = qrunner.New(config.DbDsn)
		if err != nil {
			log.Println("unable to connect to db.", err)
		}
		g.GET("/connect", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "<div class='notification is-warning'>Using database connection specified from environment variable</div>")
		})
	} else {
		connect := connect{ctr}
		g.GET("/connect", connect.Get)
		g.POST("/connect", connect.Post)
	}

	g.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong\n")
	})
}
