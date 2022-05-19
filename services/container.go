package services

import (
	"log"
  "fmt"

	"sql-ui/services/qrunner"
	"sql-ui/config"

	"github.com/gin-gonic/gin"
	_ "github.com/xo/usql/drivers/sqlite3"
)


type Container struct {
  Qrunner *qrunner.Qrunner

	// Web stores the web framework
	Web *gin.Engine

	// Config stores the application configuration
	Config *config.Config

	// TemplateRenderer stores a service to easily render and cache templates
	TemplateRenderer *TemplateRenderer
}


func NewContainer() *Container {
	c := new(Container)
	c.initWeb()
  c.initConfig()
	c.initTemplateRenderer()

  // TODO remove sqlite import 
  qr, err := qrunner.New("sqlite:///tmp/tt/db.db")
  c.Qrunner = qr
  if err != nil {
    log.Fatal(err)
  }

  return c 
}

// Shutdown shuts the Container down and disconnects all connections
func (c *Container) Shutdown() error {
  return c.Qrunner.Close()
}


// initWeb initializes the web framework
func (c *Container) initWeb() {
	c.Web = gin.Default()
}

func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c.Config = &cfg
}

func (c *Container) initTemplateRenderer() {
	c.TemplateRenderer = NewTemplateRenderer(c.Config)
}
