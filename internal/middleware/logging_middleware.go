package middleware

import (
	"fmt"
	"log/slog"
	"social-network-service/internal/log"

	"github.com/gin-gonic/gin"
)

func NewLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		attrs := []slog.Attr{}

		ctx := c.Request.Context()

		requestId, ok := ctx.Value("x-request-id").(string)

		if ok {
			attrs = append(attrs, slog.String("x-request-id", requestId))
		}

		attrs = append(attrs, slog.String("endpoint", c.Request.URL.Path))

		ctx = log.AddToContext(ctx, attrs)

		c.Request = c.Request.WithContext(ctx)

		slog.InfoContext(ctx, fmt.Sprintf("%s endpoint called", c.Request.URL.Path))

		c.Next()

		for _, e := range c.Errors {
			slog.ErrorContext(ctx, fmt.Sprintf("Got an error: %s", e.Error()))
		}

		slog.InfoContext(
			ctx,
			fmt.Sprintf(
				"%s endpoint handling finished with status %v",
				c.Request.URL.Path,
				c.Writer.Status(),
			),
		)
	}
}
