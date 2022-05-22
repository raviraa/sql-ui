package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	toml "github.com/pelletier/go-toml"
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
	// Timeout stores default http server timeout
	HttpTimeout = time.Second * 180
	// CacheControl stores default cache time for StaticDir
	CacheControl = time.Hour * 24
)

var AppDebug = false

type HistEntry struct {
	Query    string
	LastRun  time.Time
	RunCount int
}

// Config stores complete configuration
type Config struct {
	PagerSize    int         `comment:"Number of rows in each page when browsing tables"`
	HistEntries  []HistEntry `comment:"History of queries run"`
	HistEntryMax int         `comment:"Maximum number of HistEntries to keep"`

	// histEntriesMap keeps track of HistEntries in memory by Query
	histEntriesMap map[string]HistEntry
}

// TODO toml config, history

// GetConfig loads and returns configuration
func GetConfig() (Config, error) {
	var cfg Config = Config{
		PagerSize:      20,
		HistEntryMax:   9,
		histEntriesMap: make(map[string]HistEntry),
	}
	err := cfg.readConf()
	if err != nil {
		log.Println("Failed to read config from: ", confLocation())
	}
	return cfg, nil
}

func (c *Config) readConf() error {
	b, err := ioutil.ReadFile(confLocation())
	if err != nil {
		log.Println("Failure in reading config ", err)
		return err
	}
	if err = toml.Unmarshal(b, c); err != nil {
		log.Println("Failure in parsing config ", err)
		return err
	}
	c.makeHistEntryMap()
	return nil
}

func confLocation() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		log.Println("configdir error: ", err)
		confdir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	return path.Join(confdir, "sql-ui.toml")
}

func (c *Config) SaveConf() error {
	f, err := os.OpenFile(confLocation(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	enc := toml.NewEncoder(f)

	return enc.Encode(*c)
}

func (c *Config) AddHistEntry(query string) {
}

func (c *Config) makeHistEntryMap() {
	// for _, he := range c.HistEntries {

}
