package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)


type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService} //! passing userService
}

func(h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas passing sebagai parameter service

	var input user.RegisterUserInput
	
	err := c.ShouldBindJSON(&input) //! validasi di lakukan di sini

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage) //! entity tidak bisa di proses 

		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, _ := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.ApiResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "tokentokentoken")

	response := helper.ApiResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}


func(h *userHandler) Login(c *gin.Context) {
	/*
		user memasukkan input
		input di tangkap handler
		mapping dari input user ke input struct
		input struct passing ke service
		di service mencari dgn bantuan repository user dengan email x
	*/

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Login account failed", http.StatusUnprocessableEntity, "error", errorMessage) //! entity tidak bisa di proses 
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ApiResponse("Login account failed", http.StatusUnprocessableEntity, "error", errorMessage) //! entity tidak bisa di proses, karena format email salah, id dan email tidak di temukan dan password salah
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, "tokentokentoken")

	response := helper.ApiResponse("Login has been Successfully", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}


func(h *userHandler) CheckEmailHasBeenRegister(c *gin.Context) {
	/*
		ada input email dari user
		input email di mapping ke struct input
		struct input di passing ke service
		service akan memanggil repository, jika email sudah ada atau belum
		repository akan melakukan query ke database
	*/

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage) //! entity tidak bisa di proses 
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.ApiResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage) //! entity tidak bisa di proses 
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.ApiResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}