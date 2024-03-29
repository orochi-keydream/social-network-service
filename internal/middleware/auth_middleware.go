package middleware

import (
	"social-network-service/internal/model"

	"github.com/gin-gonic/gin"
)

type JwtService interface {
	CheckAccess(c *gin.Context) (model.UserId, error)
}

func NewAuthMiddleware(jwtService JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwtService.CheckAccess(c)

		if err != nil {
			c.Error(err)
			return
		}

		c.Next()
	}
}
