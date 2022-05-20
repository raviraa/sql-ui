package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)



func CacheControl(maxAge time.Duration) gin.HandlerFunc {
  return func(ctx *gin.Context) {
			v := "no-cache, no-store"
			if maxAge > 0 {
				v = fmt.Sprintf("public, max-age=%.0f", maxAge.Seconds())
			}
    ctx.Header("Cache-Control", v)
    ctx.Next()
  }
}
