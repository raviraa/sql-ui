//go:build sampledb
// sampledb embeds chinook sample sqlite db, extracts and writs the db to config.SampleDbPath.
// Use tag sampledb to enable and build this. Nothing is included without the tag.
package sampledb

import (
	"compress/gzip"
	"embed"
	_ "embed"
	"io"
	"log"
	"os"

	"github.com/raviraa/sql-ui/config"
)

//go:embed sample.db.gz
var gzfile embed.FS

func init() {
  log.Println("Initializing sampledb to: ", config.SampleDbPath)

	gz, err := gzfile.Open("sample.db.gz")
	if err != nil {
		log.Println(err)
		return
	}
	zr, err := gzip.NewReader(gz)
	if err != nil {
		log.Println(err)
		return
	}
	wf, err := os.OpenFile(config.SampleDbPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := io.Copy(wf, zr); err != nil {
		log.Println(err)
		return
	}
	if err := zr.Close(); err != nil {
		log.Println(err)
		return
	}
	if err := wf.Close(); err != nil {
		log.Println(err)
		return
	}
}
