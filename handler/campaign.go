package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
	tangkap parameter di handler
	handler ke service
	service menentukan repository mana yang akan di call
	repository akses ke db : FindALl, FindUserByID
*/

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func(h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.ApiResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	
	c.JSON(http.StatusOK, response)
}

func(h *campaignHandler) GetCampaign(c *gin.Context) {
	/*
		memasukkan request user berupa parameter id nya
		handler: mapping id yg di url ke struct input => service, memanggil formatter untuk melakukan formatting
		service: input struct input untuk menangkap id di url, passsing ke service
		repository: get campaign by id
	*/

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to get detail campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.ApiResponse("Failed to get detail campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	} 

	response := helper.ApiResponse("Success get Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func(h *campaignHandler) CreateCampaign(c *gin.Context) {
	/*
		tangkap parameter dari user ke input struct
		ambil current user dari jwt/handler
		panggil service , parameter input struct (membuat slug => berdasarkan nama campaign)
		panggil repository untuk simpan data campaign baru
	*/

	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.ApiResponse("Failed Create Campaign", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.ApiResponse("Failed Create Campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	response := helper.ApiResponse("Success Create Campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))

	c.JSON(http.StatusOK, response)
}

func(h *campaignHandler) UpdateCampaign(c *gin.Context) {
	/*
		User memasukkan input
		handler menangkap inputan dari user
		mapping dari input ke input struct (adri user dan uri)
		passing ke service
		Service memanggil atau menggunakan function yang ada di repository
		repository update data campaign
	*/

	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.ApiResponse("Failed Update Campaign", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.ApiResponse("Failed Update Campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	response := helper.ApiResponse("Success Update Campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))

	c.JSON(http.StatusOK, response)
}

/*
	! setiap membuat sebuah handler perlu di daftarkan routingnya
*/