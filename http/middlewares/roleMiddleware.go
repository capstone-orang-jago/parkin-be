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

		userRole, ok := role.(domain.UserRole)
		if !ok {
			helpers.Response(c, domain.HttpResponse{
				Code:    http.StatusForbidden,
				Message: "invalid role",
			})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
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
