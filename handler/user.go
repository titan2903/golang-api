package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService} //! passing userService
}

func(h *userHandler) RegisterUser(c *gin.Context) {
	/*
		tangkap input dari user
	 	map input dari user ke struct RegisterUserInput
	 	struct di atas passing sebagai parameter service
	*/

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

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.ApiResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

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

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.ApiResponse("Login account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)

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

func(h *userHandler) UploadAvatar(c *gin.Context) {
	/*
		inputan berupa http form jadi tidak perlu membuat struct
		input dari user
		simpan gambar di folder "images/"
		di service panggil repository
		JWT (sementara hardcode, seakan2 user yang login ID = 1)
		repository ambil data user yang ID = 1
		repository update data user simpan lokasi file
	*/

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
	}

	//! harusnya dapat dari JWT, tapi belum!!
	userID := 2
	// path := "images/" + file.Filename //! lokasi file di simpan
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename) //! agar file name unique mengikudi id user

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}