package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.NewString()

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "x-request-id", requestId)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
