package middlewares

import (
	"strings"

	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.Response(c, domain.HttpResponse{
				Code:    401,
				Message: "authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.Response(c, domain.HttpResponse{
				Code:    401,
				Message: "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := helpers.ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			helpers.Response(c, domain.HttpResponse{
				Code:    401,
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		userID, ok := token["sub"].(string)
		if !ok {
			helpers.Response(c, domain.HttpResponse{
				Code:    401,
				Message: "invalid token claims",
			})
			c.Abort()
			return
		}

		c.Set("sub", userID)
		c.Next()
	}
}
