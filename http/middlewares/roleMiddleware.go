package middlewares

import (
	"net/http"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			helpers.Response(c, domain.HttpResponse{
				Code:    http.StatusForbidden,
				Message: "role not found",
			})
			c.Abort()
			return
		}

		var userRoleStr string
		switch v := role.(type) {
		case string:
			userRoleStr = v
		case domain.UserRole:
			userRoleStr = string(v)
		default:
			helpers.Response(c, domain.HttpResponse{
				Code:    http.StatusForbidden,
				Message: "invalid role type in context",
			})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRoleStr == string(allowed) {
				c.Next()
				return
			}
		}

		helpers.Response(c, domain.HttpResponse{
			Code:    http.StatusForbidden,
			Message: "forbidden",
		})
		c.Abort()
	}
}
