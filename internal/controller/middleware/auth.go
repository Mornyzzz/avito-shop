package middleware

import (
	"avito-shop/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})

		jwt.AuthMiddleware(next).ServeHTTP(c.Writer, c.Request)

		if c.Writer.Written() {
			c.Abort()
		}
	}
}
