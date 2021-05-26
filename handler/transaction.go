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


func(h *transactionHandler) GetUserTransactions(c *gin.Context) {
/*
	GetUserTransaction
	handler: 
		- ambil nilai user dari jwt atau middleware
	service
	repository:
		- ambil data transaction (preload data campaign)
*/

	currentUser := c.MustGet("currentUser").(user.User) //! get id ddari user yg login melalui jwt
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionByUserID(userID)
	if err != nil {
		response := helper.ApiResponse("Failed get users's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	response :=  helper.ApiResponse("Success get campaign's transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func(h *transactionHandler) CreateTransaction(c *gin.Context) {
	/*
		ada input dari user memasukkan jumlah uang yang di input
		handler menangkap input dan kemudian di mapping ke input struct nya
		di dalam handled memanggil service buat transaksi, akan di record datanya ke database, manggil sistem midtrans untuk mendaftarkan transaksinya. balikannya berupo token.
		service memanggil repository create new transaction data
	*/

	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Failed Create transaction", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.ApiResponse("Failed Create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	response := helper.ApiResponse("Success Create transaction", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))

	c.JSON(http.StatusOK, response)
}


func(h *transactionHandler) GetNotification(c *gin.Context) { //! yang mengakses endpoint ini bukan client, tetapi yang mengaksesnya yaitu midtrans
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	err = h.service.PaymentProcess(input)
	if err != nil {
		response := helper.ApiResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return;
	}

	c.JSON(http.StatusOK, input)
}