package middleware

import (
	"errors"
	"log"
	"net/http"
	"social-network-service/internal/model"

	"github.com/gin-gonic/gin"
)

func NewErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors[0]

		var a *model.ClientError

		if errors.As(err, &a) {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		var b *model.UnauthenticatedError

		if errors.As(err, &b) {
			log.Println(err)
			c.Status(http.StatusUnauthorized)
			return
		}

		var d *model.ForbiddenError

		if errors.As(err, &d) {
			log.Println(err)
			c.Status(http.StatusForbidden)
			return
		}

		var f *model.NotFoundError

		if errors.As(err, &f) {
			log.Println(err)
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusInternalServerError)
	}
}
