package handler

import (
	"golang-api-crowdfunding/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type sessionHanlder struct {
	userService user.Service
}

func NewSessionHandler(userService user.Service) *sessionHanlder {
	return &sessionHanlder{userService}
}

func (h *sessionHanlder) FormLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "session_new.html", nil)
}

func (h *sessionHanlder) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user, err := h.userService.Login(input)
	if err != nil || user.Role != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Set("userName", user.Name)
	session.Save() //! simpan data session

	c.Redirect(http.StatusFound, "/users")
}

func (h *sessionHanlder) Destroy(c *gin.Context) {
	sessions := sessions.Default(c)
	sessions.Clear() //! menghapus data session
	sessions.Save()
	c.Redirect(http.StatusFound, "/login")
}
