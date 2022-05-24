package htmx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//  Headers (https://htmx.org/docs/#requests)
const (
	HeaderRequest            = "HX-Request"
	HeaderBoosted            = "HX-Boosted"
	HeaderTrigger            = "HX-Trigger"
	HeaderTriggerName        = "HX-Trigger-Name"
	HeaderTriggerAfterSwap   = "HX-Trigger-After-Swap"
	HeaderTriggerAfterSettle = "HX-Trigger-After-Settle"
	HeaderTarget             = "HX-Target"
	HeaderPrompt             = "HX-Prompt"
	HeaderPush               = "HX-Push"
	HeaderRedirect           = "HX-Redirect"
	HeaderRefresh            = "HX-Refresh"
)

type (
	// Request contains data that HTMX provides during requests
	Request struct {
		Enabled     bool
		Boosted     bool
		Trigger     string
		TriggerName string
		Target      string
		Prompt      string
	}

	// Response contain data that the server can communicate back to HTMX
	Response struct {
		Push               string
		Redirect           string
		Refresh            bool
		Trigger            string
		TriggerAfterSwap   string
		TriggerAfterSettle string
		NoContent          bool
	}
)

// GetRequest extracts HTMX data from the request
func GetRequest(ctx *gin.Context) Request {
	return Request{
		Enabled:     ctx.Request.Header.Get(HeaderRequest) == "true",
		Boosted:     ctx.Request.Header.Get(HeaderBoosted) == "true",
		Trigger:     ctx.Request.Header.Get(HeaderTrigger),
		TriggerName: ctx.Request.Header.Get(HeaderTriggerName),
		Target:      ctx.Request.Header.Get(HeaderTarget),
		Prompt:      ctx.Request.Header.Get(HeaderPrompt),
	}
}

// Apply applies data from a Response to a server response
func (r Response) Apply(ctx *gin.Context) {
	if r.Push != "" {
		ctx.Header(HeaderPush, r.Push)
	}
	if r.Redirect != "" {
		ctx.Header(HeaderRedirect, r.Redirect)
	}
	if r.Refresh {
		ctx.Header(HeaderRefresh, "true")
	}
	if r.Trigger != "" {
		ctx.Header(HeaderTrigger, r.Trigger)
	}
	if r.TriggerAfterSwap != "" {
		ctx.Header(HeaderTriggerAfterSwap, r.TriggerAfterSwap)
	}
	if r.TriggerAfterSettle != "" {
		ctx.Header(HeaderTriggerAfterSettle, r.TriggerAfterSettle)
	}
	if r.NoContent {
		// ctx.Response().Status = http.StatusNoContent
		ctx.Status(http.StatusNoContent)
	}
}
