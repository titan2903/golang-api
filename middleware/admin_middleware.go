package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions := sessions.Default(c)

		userIDSession := sessions.Get("userID")
		if userIDSession == nil { //! user dalam kondisi tidak login
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}
}