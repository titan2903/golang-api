package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/middleware"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:user1234@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//! Auth
	authService := auth.NewService()

	//! Users
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, authService)

	//!Campaigns
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	//! Payment
	paymentService := payment.NewService()

	// ! Transaction
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	
	router := gin.Default()
	api := router.Group("/api/v1")

	//! Router access immage
	router.Static("/images", "./images")

	//! Router Users
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailHasBeenRegister)
	api.POST("/avatars", middleware.AuthMiddleware(authService, userService) ,userHandler.UploadAvatar)
	api.GET("/users/fetch", middleware.AuthMiddleware(authService, userService), userHandler.FetchUser)

	//!Router Campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("campaign-image", middleware.AuthMiddleware(authService, userService), campaignHandler.UploadCampaignImage)

	//!Router Handler
	api.GET("/campaigns/:id/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", middleware.AuthMiddleware(authService,userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.Use(cors.Default()) // ! Allow cors
	router.Run() //! default PORT 8080

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Tes simpan dari service"
	// userInput.Email = "test@gmail.com"
	// userInput.Occupation = "anak band"
	// userInput.Password = "user1234"
	// userService.RegisterUser(userInput)

	// var users []user.User
	// db.Find(&users) //! type harus pointer

	// for _, user := range users {
	// 	fmt.Println(user.Name)
	// 	fmt.Println(user.Email)
	// }

	// router := gin.Default() //! daftarkan router default dari Gin
	// router.GET("/handler", handler) //! memanggil handler functionnya
	// router.Run() //! menjalankan router
}

//! input (memasukkan data atau mengirim request dari client) -> Handler (mapping input ke struct) -> memanggil Service (melakukan bisnis proses, mapping struct) -> repository(akses ke database, berupa CRUD) -> memanggil DB