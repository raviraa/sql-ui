package routes

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple example of how to test routes and their markup using the test HTTP server spun up within
// this test package
func TestPing_Get(t *testing.T) {
	doc := request(t).
		setRoute("ping").
		get().
		assertStatusCode(http.StatusOK).
		toDoc()

	assert.Equal(t, "pong", doc.Text())

	// h1 := doc.Find("h1.title")
	// assert.Len(t, h1.Nodes, 1)
	// assert.Equal(t, "About", h1.Text())
}
