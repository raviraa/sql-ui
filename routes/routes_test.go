package routes

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/raviraa/sql-ui/config"
	"github.com/raviraa/sql-ui/services"
	"github.com/raviraa/sql-ui/services/qrunner"
	_ "github.com/xo/usql/drivers/sqlite3"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var (
	srv *httptest.Server
	c   *services.Container
)

func TestMain(m *testing.M) {
	config.AppEnvironment = config.EnvTest
	// Start a new container
	c = services.NewContainer()
	var err error
	c.Qrunner, err = qrunner.New("../services/qrunner/testdata/data.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// Start a test HTTP server
	BuildRouter(c)
	srv = httptest.NewServer(c.Web)

	// Run tests
	exitVal := m.Run()

	// Shutdown the container and test server
	if err := c.Shutdown(); err != nil {
		panic(err)
	}
	srv.Close()

	os.Exit(exitVal)
}

type httpRequest struct {
	route   string
	client  http.Client
	headers map[string]string
	body    url.Values
	t       *testing.T
}

func request(t *testing.T) *httpRequest {
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	r := httpRequest{
		t:       t,
		body:    url.Values{},
		headers: make(map[string]string),
		client: http.Client{
			Jar: jar,
		},
	}
	return &r
}

func (h *httpRequest) setClient(client http.Client) *httpRequest {
	h.client = client
	return h
}

func (h *httpRequest) setRoute(route string, params ...interface{}) *httpRequest {
	// h.route = srv.URL + c.Web.Reverse(route, params)
	h.route = srv.URL + "/" + route
	return h
}

func (h *httpRequest) setHeader(key, val string) *httpRequest {
	h.headers[key] = val
	return h
}

func (h *httpRequest) addBody(key, val string) *httpRequest {
	h.headers["Content-Type"] = "application/x-www-form-urlencoded"
	h.body.Add(key, val)
	return h
}

func (h *httpRequest) get() *httpResponse {
	// resp, err := h.client.Get(h.route)
	req, err := http.NewRequest("GET", h.route, nil)
	require.NoError(h.t, err)
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	resp, err := h.client.Do(req)
	require.NoError(h.t, err)
	r := httpResponse{
		t:        h.t,
		Response: resp,
	}
	return &r
}

func (h *httpRequest) post() *httpResponse {

	// Make the POST requests

	// return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	// resp, err := h.client.PostForm(h.route, h.body)
	req, err := http.NewRequest("POST", h.route, strings.NewReader(h.body.Encode()))
	require.NoError(h.t, err)
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	resp, err := h.client.Do(req)
	require.NoError(h.t, err)
	r := httpResponse{
		t:        h.t,
		Response: resp,
	}
	return &r
}

type httpResponse struct {
	*http.Response
	t *testing.T
}

func (h *httpResponse) assertStatusCode(code int) *httpResponse {
	assert.Equal(h.t, code, h.Response.StatusCode)
	return h
}

func (h *httpResponse) assertRedirect(t *testing.T, route string, params ...interface{}) *httpResponse {
	assert.Equal(t, route, h.Header.Get("Location"))
	return h
}

func (h *httpResponse) toDoc() *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(h.Body)
	require.NoError(h.t, err)
	err = h.Body.Close()
	assert.NoError(h.t, err)
	return doc
}
