package routes

import (
	"net/http"
	"sql-ui/config"
	"sql-ui/controller"
	"sql-ui/middleware"
	"sql-ui/services"

	"github.com/gin-gonic/gin"
)

func BuildRouter(c *services.Container) {

	// c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
	// 	Static(config.StaticPrefix, config.StaticDir)
	// TODO staticdir from embedfs
	c.Web.Group("").
		Use(middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
		Static(config.StaticPrefix, config.StaticDir)

	g := c.Web.Group("")
	g.Use(
		gin.ErrorLogger(),
		gin.Recovery(),
	)

	// Base controller
	ctr := controller.NewController(c)
	// TODO controller
	userRoutes(c, g, ctr)

}

// func navRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
func userRoutes(c *services.Container, g *gin.RouterGroup, ctr controller.Controller) {

	home := home{Controller: ctr}
	g.GET("/", home.Get)

	query := query{Controller: ctr}
	g.GET("/query", query.Get)
	g.POST("/query", query.Post)

	tables := tables{ctr}
	g.GET("/tables", tables.Get)
	g.POST("/tables", tables.Post)

	g.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
}
