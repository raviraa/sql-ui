package routes

import (
	"net/http"
	"testing"

	"sql-ui/htmx"
	"github.com/stretchr/testify/assert"
)

func TestQuery_Get(t *testing.T) {
	doc := request(t).
		setRoute("").
		get().
		assertStatusCode(http.StatusOK).
		toDoc()

	title := doc.Find("title")
  assert.Equal(t, "Query", title.Text())
  assert.Len(t, title.Nodes, 1)
  txarea := doc.Find("textarea[name=query]")
  assert.Len(t, txarea.Nodes, 1)
  assert.Equal(t, "/query", txarea.AttrOr("hx-post", ""))
}

func TestQuery_Post(t *testing.T) {
	doc := request(t).
		setRoute("query").
    setHeader(htmx.HeaderRequest, "true").
    addBody("query", "select * from t1 order by val limit 3").
		post().
		assertStatusCode(http.StatusOK).
		toDoc()

	// title := doc.Find("title")
  // assert.Equal(t, "Query", title.Text())
  tds := doc.Find("td")
  assert.Len(t, tds.Nodes, 6)
  assert.Equal(t, "11one22two33three", tds.Text())
}

func TestQueryErrPost(t *testing.T) {
	doc := request(t).
		setRoute("query").
    setHeader(htmx.HeaderRequest, "true").
    addBody("query", "select * from tttt").
		post().
		assertStatusCode(http.StatusOK).
		toDoc()

  tds := doc.Find("td")
  assert.Len(t, tds.Nodes, 0)
  assert.Contains(t, doc.Text(), "sqlite3: 1: no such table: tttt")
}
