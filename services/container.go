package services

import (
	"fmt"
	"log"

	"sql-ui/config"
	"sql-ui/services/qrunner"

	"github.com/gin-gonic/gin"
	// _ "github.com/xo/usql/drivers/sqlite3"
	_ "sql-ui/internal"
	// _ "github.com/xo/usql/internal"
)

// Container keep the state of application
type Container struct {
	Qrunner *qrunner.Qrunner

	// Web stores the web framework
	Web *gin.Engine

	// Config stores the application configuration
	Config *config.Config

	// TemplateRenderer stores a service to easily render and cache templates
	TemplateRenderer *TemplateRenderer

	// Query is sql entered in textarea
	Query string
}

// NewContainer initializes web, templates etc.
func NewContainer() *Container {
	log.SetFlags(log.Lshortfile)
	c := new(Container)
	c.initWeb()
	c.initConfig()
	c.initTemplateRenderer()
  c.Query = c.Config.GetHistEntryRecent()

	return c
}

// Shutdown shuts the Container down and disconnects all connections
func (c *Container) Shutdown() error {
	if config.AppEnvironment != config.EnvTest {
		log.Println("Stopping server, and saving conf")
		if err := c.Config.SaveConf(); err != nil {
			log.Println("error saving conf: ", err)
		}
	}
	return c.Qrunner.Close()
}

// initWeb initializes the web framework
func (c *Container) initWeb() {
	if config.AppEnvironment != config.EnvDev {
		gin.SetMode(gin.ReleaseMode)
	}
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
