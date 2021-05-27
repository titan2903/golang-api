package handler

import (
	"bwastartup/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func(h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		//! DO Handle
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return;
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users}) //! datanya berupa key value (map)
}

func(h *userHandler) FormCreateUser(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", nil)
}

func(h *userHandler) CreateUser(c *gin.Context) {
	var input user.FormCreateUserInput

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_new.html", input) //! menampilkan lagi data yang di input user jika error
		return;
	}

	registerInput := user.RegisterUserInput{}
	registerInput.Name = input.Name
	registerInput.Email = input.Email
	registerInput.Occupation = input.Occupation
	registerInput.Password = input.Password

	_, err = h.userService.RegisterUser(registerInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return;
	}

	c.Redirect(http.StatusFound, "/users")
}

func(h *userHandler) FormUpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	registerUser, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := user.FormUpdateUserInput{}
	input.ID = registerUser.ID
	input.Name = registerUser.Name
	input.Email = registerUser.Email
	input.Occupation = registerUser.Occupation

	c.HTML(http.StatusOK, "user_edit.html", input)
}

func(h *userHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input user.FormUpdateUserInput

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		c.HTML(http.StatusOK, "user_edit.html", input) //! menampilkan lagi data yang di input user jika error
		return;
	}

	input.ID = id
	_, err = h.userService.UpdateUser(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}