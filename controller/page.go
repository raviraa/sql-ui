package controller

import (
	"net/http"

	"sql-ui/htmx"
	"sql-ui/services"

	"github.com/gin-gonic/gin"
	// "github.com/mikestefanello/pagoda/msg"
)

// Page consists of all data that will be used to render a page response for a given controller.
// While it's not required for a controller to render a Page on a route, this is the common data
// object that will be passed to the templates, making it easy for all controllers to share
// functionality both on the back and frontend. The Page can be expanded to include anything else
// your app wants to support.
// Methods on this page also then become available in the templates, which can be more useful than
// the funcmap if your methods require data stored in the page, such as the context.
type Page struct {

	// Title stores the title of the page
	Title string

	// Context stores the request context
	Context *gin.Context

	// Path stores the path of the current request
	Path string

	// URL stores the URL of the current request
	URL string

	// Data stores whatever additional data that needs to be passed to the templates.
	// This is what the controller uses to pass the content of the page.
	Data interface{}

	// Layout stores the name of the layout base template file which will be used when the page is rendered.
	// This should match a template file located within the layouts directory inside the templates directory.
	// The template extension should not be included in this value.
	Layout string

	// Name stores the name of the page as well as the name of the template file which will be used to render
	// the content portion of the layout template.
	// This should match a template file located within the pages directory inside the templates directory.
	// The template extension should not be included in this value.
	Name string

	// HtmxBase is is base layout used in case of htmx requests
	HtmxBase string

	// TemplateExtra stores extra template names in pages directory that are needed to execute the page
	TemplateExtra []string

	// StatusCode stores the HTTP status code that will be returned
	StatusCode int

	// Headers stores a list of HTTP headers and values to be set on the response
	Headers map[string]string

	HTMX struct {
		Request  htmx.Request
		Response *htmx.Response
	}

	Container *services.Container
}

// NewPage creates and initiatizes a new Page for a given request context
func NewPage(ctx *gin.Context, container *services.Container) Page {
	p := Page{
		Context:    ctx,
		Container:  container,
		HtmxBase:   "htmx",
		Path:       ctx.Request.URL.Path,
		URL:        ctx.Request.RequestURI,
		StatusCode: http.StatusOK,
		Headers:    make(map[string]string),
	}

	p.HTMX.Request = htmx.GetRequest(ctx)

	return p
}
