package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)


type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

/*
	parameter di uri (id dari ca,paign)
	tangkap parameter mapping input struct ShouldBindUri()
	memanggil service atau struct di passing ke service, input struct sebagai parameter
	service, berbekal campaign id bisa memanggil repository
	repository mencari data transaction suatu campaign
*/


func(h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	currentUser := c.MustGet("currentUser").(user.User) //! melakukan auth user, hanya user yang memiliki item tsb bisa melakukabn update
	input.User = currentUser

	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.ApiResponse("Failed get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	response := helper.ApiResponse("Success get campaign's transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))

	c.JSON(http.StatusOK,response)
}