package controller

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"

	// "sql-ui/context"
	// "sql-ui/htmx"
	"sql-ui/services"

	"github.com/gin-gonic/gin"
	// "github.com/labstack/echo/v4"
)

// Controller provides base functionality and dependencies to routes.
// The proposed pattern is to embed a Controller in each individual route struct and to use
// the router to inject the container so your routes have access to the services within the container
type Controller struct {
	// Container stores a services container which contains dependencies
	Container *services.Container
}

// NewController creates a new Controller
func NewController(c *services.Container) Controller {
	return Controller{
		Container: c,
	}
}

// RenderPage renders a Page as an HTTP response
func (c *Controller) RenderPage(ctx gin.Context, page Page) error {
	var buf *bytes.Buffer
	var err error

	// Page name is required
	if page.Name == "" {
		log.Println("page render failed due to missing name")
		// return echo.NewHTTPError(http.StatusInternalServerError)
		return errors.New("missing page name")
	}

	// Check if this is an HTMX non-boosted request which indicates that only partial
	// content should be rendered
	if page.HTMX.Request.Enabled && !page.HTMX.Request.Boosted {
		// Parse and execute the templates only for the content portion of the page
		// The templates used for this partial request will be:
		// 1. The base htmx template which omits the layout and only includes the content template
		// 2. The content template specified in Page.Name
		// 3. All templates within the components directory
		// Also included is the function map provided by the funcmap package
		buf, err = c.Container.TemplateRenderer.
			Parse().
			Group("page:htmx").
			Key(page.Name).
			Base("htmx").
			Files(
				"htmx",
				fmt.Sprintf("pages/%s", page.Name),
			).
			Directories("components").
			Execute(page)
	} else {
		// Parse and execute the templates for the Page
		// As mentioned in the documentation for the Page struct, the templates used for the page will be:
		// 1. The layout/base template specified in Page.Layout
		// 2. The content template specified in Page.Name
		// 3. All templates within the components directory
		// Also included is the function map provided by the funcmap package
		buf, err = c.Container.TemplateRenderer.
			Parse().
			Group("page").
			Key(page.Name).
			Base(page.Layout).
			Files(
				fmt.Sprintf("layouts/%s", page.Layout),
				fmt.Sprintf("pages/%s", page.Name),
			).
			Directories("components").
			Execute(page)
	}

	if err != nil {
		log.Println("failed to parse and execute templates:", err)
		// return echo.NewHTTPError(http.StatusInternalServerError)
		ctx.String(http.StatusInternalServerError, "template error")
		return errors.New("template error")
	}

	// Set the status code
	// ctx.Response().Status = page.StatusCode
	ctx.Status(page.StatusCode)

	// Set any headers
	for k, v := range page.Headers {
		// ctx.Response().Header().Set(k, v)
		ctx.Header(k, v)
	}

	// Apply the HTMX response, if one
	if page.HTMX.Response != nil {
		page.HTMX.Response.Apply(ctx)
	}

	// Cache this page, if caching was enabled
	// c.cachePage(ctx, page, buf)

	// return ctx.HTMLBlob(ctx.Response().Status, buf.Bytes())
	_, err = ctx.Writer.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
	}
	return err
}

/*
// Redirect redirects to a given route name with optional route parameters
func (c *Controller) Redirect(ctx echo.Context, route string, routeParams ...interface{}) error {
	url := ctx.Echo().Reverse(route, routeParams)

	if htmx.GetRequest(ctx).Boosted {
		htmx.Response{
			Redirect: url,
		}.Apply(ctx)

		return nil
	} else {
		return ctx.Redirect(http.StatusFound, url)
	}
}

// Fail is a helper to fail a request by returning a 500 error and logging the error
func (c *Controller) Fail(ctx echo.Context, err error, log string) error {
	if context.IsCanceledError(err) {
		return nil
	}
	ctx.Logger().Errorf("%s: %v", log, err)
	return echo.NewHTTPError(http.StatusInternalServerError)
}
*/
