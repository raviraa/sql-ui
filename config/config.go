package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"time"

	toml "github.com/pelletier/go-toml"
)

const (
	// TemplateDir stores the name of the directory that contains templates
	TemplateDir = "templates"
	// TemplateExt stores the extension used for the template files
	TemplateExt = ".html"
	// StaticPrefix stores the URL prefix used when serving static files
	StaticPrefix = "files"
	// Timeout stores default http server timeout
	HttpTimeout = time.Second * 180
	// CacheControl stores default cache time for StaticDir
	CacheControl = time.Hour * 24
)

type Environment string

var (
	EnvTest Environment = "test"
	EnvDev  Environment = "dev"
)

var AppEnvironment = EnvDev

var SampleDbPath = "/tmp/sample.db"

// DbDsn forces db to connect from env var DB_DSN. Will not be able to connect to other databases
var DbDsn = ""

type HistEntry struct {
	Query    string
	LastRun  time.Time
	RunCount int
}

// Config stores complete configuration
type Config struct {
	PagerSize        int         `comment:"Number of rows in each page when browsing tables"`
	HistEntries      []HistEntry `comment:"History of queries run"`
	HistEntryMax     int         `comment:"Maximum number of HistEntries to keep"`
	ConnectDSN       []string    `comment:"History of connected database server dsn"`
	ConnectDSNMax    int         `comment:"Maximum number of DSN entries to keep"`
	OpenInWebBrowser bool        `comment:"Open server url in web browser on startup"`

	// histEntriesMap keeps track of HistEntries in memory by Query
	histEntriesMap map[string]*HistEntry
}

// GetConfig loads and returns configuration
func GetConfig() (Config, error) {
	var cfg Config = Config{
		PagerSize:        20,
		HistEntryMax:     99,
		ConnectDSNMax:    99,
		histEntriesMap:   make(map[string]*HistEntry),
		OpenInWebBrowser: true,
	}
	err := cfg.readConf()
	if err != nil {
		log.Println("Failed to read config from: ", confLocation())
	}
	if dbdsn := os.Getenv("DB_DSN"); dbdsn != "" {
		DbDsn = dbdsn
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

var confLocationStr string

func confLocation() string {
	if confLocationStr != "" {
		return confLocationStr
	}
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
	c.syncHistEntryMap()
	enc := toml.NewEncoder(f)
	return enc.Encode(*c)
}

func (c *Config) AddHistEntry(query string) {
	he, ok := c.histEntriesMap[query]
	if ok {
		he.RunCount += 1
		he.LastRun = time.Now()
	} else {
		c.histEntriesMap[query] = &HistEntry{
			Query:    query,
			LastRun:  time.Now(),
			RunCount: 1,
		}
	}
}

func (c *Config) GetHistEntryRecent() string {
	hentries := c.GetHistEntries()
	if len(hentries) > 0 {
		return hentries[0].Query
	}
	return ""
}

func (c *Config) GetHistEntries() []HistEntry {
	c.syncHistEntryMap()
	hentries := c.HistEntries
	sort.Slice(hentries, func(i, j int) bool {
		if hentries[i].RunCount == hentries[j].RunCount {
			return hentries[i].LastRun.Before(hentries[j].LastRun)
		}
		return hentries[i].RunCount < hentries[j].RunCount
	})
	return hentries
}

func (c *Config) makeHistEntryMap() {
	for idx := range c.HistEntries {
		c.histEntriesMap[c.HistEntries[idx].Query] = &c.HistEntries[idx]
	}
}

func (h HistEntry) String() string {
	return fmt.Sprintf("{%v;%v;%v}", h.RunCount, h.LastRun.UnixMicro(), h.Query)
}

func (c *Config) syncHistEntryMap() {
	hentries := make([]HistEntry, 0)
	for _, he := range c.histEntriesMap {
		hentries = append(hentries, *he)
	}
	sort.Slice(hentries, func(i, j int) bool {
		if hentries[i].RunCount == hentries[j].RunCount {
			return hentries[i].LastRun.Before(hentries[j].LastRun)
		}
		return hentries[i].RunCount < hentries[j].RunCount
	})
	if len(hentries) > c.HistEntryMax {
		hentries = hentries[len(hentries)-c.HistEntryMax:]
	}
	c.HistEntries = hentries
}

func (c *Config) AddDSN(dsn string) {
	newdsns := []string{}
	for _, histdsn := range c.ConnectDSN {
		// Remove from history and add at the end
		if dsn != histdsn {
			newdsns = append(newdsns, histdsn)
		}
	}
	newdsns = append(newdsns, dsn)
	if len(newdsns) > c.ConnectDSNMax {
		newdsns = newdsns[len(newdsns)-c.ConnectDSNMax:]
	}
	c.ConnectDSN = newdsns
}

func (c *Config) DSNRecent() string {
	if len(c.ConnectDSN) == 0 {
		return ""
	}
	return c.ConnectDSN[len(c.ConnectDSN)-1]
}

func (c *Config) DSNsRecent() []string {
	dsns := []string{}
	for i := len(c.ConnectDSN) - 1; i >= 0; i-- {
		dsns = append(dsns, c.ConnectDSN[i])
	}
	return dsns
}
