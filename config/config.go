package config

import (
	"time"

	"github.com/joeshaw/envdecode"
)

const (
	// TemplateDir stores the name of the directory that contains templates
	TemplateDir = "templates"

	// TemplateExt stores the extension used for the template files
	TemplateExt = ".html"

	// StaticDir stores the name of the directory that will serve static files
	StaticDir = "static"

	// StaticPrefix stores the URL prefix used when serving static files
	StaticPrefix = "files"
)

type environment string

const (
	// EnvLocal represents the local environment
	EnvLocal environment = "local"
	// EnvProduction represents the production environment
	EnvProduction environment = "prod"
)

type (
	// Config stores complete configuration
	Config struct {
		Environment environment   `env:"APP_ENVIRONMENT,default=local"`
		Timeout     time.Duration `env:"APP_TIMEOUT,default=20s"`
		PagerSize   int           `env:"PagerSize,default=10"`
		CacheControl time.Duration `env:"CacheControl,default=4380h"`
	}
)

// GetConfig loads and returns configuration
func GetConfig() (Config, error) {
	var cfg Config
	err := envdecode.StrictDecode(&cfg)
	return cfg, err
}
