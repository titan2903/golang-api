package main

import (
	"fmt"
	"golang-api-crowdfunding/auth"
	"golang-api-crowdfunding/campaign"
	"golang-api-crowdfunding/config"
	"golang-api-crowdfunding/handler"
	"golang-api-crowdfunding/helper"
	"golang-api-crowdfunding/libraryloadtemplate"
	"golang-api-crowdfunding/middleware"
	"golang-api-crowdfunding/payment"
	"golang-api-crowdfunding/transaction"
	"golang-api-crowdfunding/user"

	webHandler "golang-api-crowdfunding/web/handler"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {

	db := config.ConnectDB()
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

	//! Web Static
	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandller := webHandler.NewCampaignHandler(campaignService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)

	//! Set Port
	port := helper.GoDotEnvVariable("PORT")
	if port == "" {
		port = "8000"
	}
	sPort := fmt.Sprintf(":%s", port)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware()) // ! Allow cors

	//!Session
	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("testbanana", cookieStore))

	//! HTML Render
	router.HTMLRender = libraryloadtemplate.LoadTemplates("./web/templates") //! mengeload tamplate yang ada di dalam folder template

	api := router.Group("/api/v1")

	//! Router access immage and web assets
	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")

	//! Router Users
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailHasBeenRegister)
	api.POST("/avatars", middleware.AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", middleware.AuthMiddleware(authService, userService), userHandler.FetchUser)

	//!Router Campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", middleware.AuthMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", middleware.AuthMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("campaign-image", middleware.AuthMiddleware(authService, userService), campaignHandler.UploadCampaignImage)

	//!Router Transactions
	api.GET("/campaigns/:id/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", middleware.AuthMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	//!Router Web Static Login
	router.GET("/login", sessionWebHandler.FormLogin)
	router.POST("/session", sessionWebHandler.Login)
	router.GET("/logout", sessionWebHandler.Destroy)

	//!Router Web Static Users
	router.GET("/users", middleware.AuthAdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", middleware.AuthAdminMiddleware(), userWebHandler.FormCreateUser)
	router.POST("/users", middleware.AuthAdminMiddleware(), userWebHandler.CreateUser)
	router.GET("/users/edit/:id", middleware.AuthAdminMiddleware(), userWebHandler.FormUpdateUser)
	router.POST("/users/update/:id", middleware.AuthAdminMiddleware(), userWebHandler.UpdateUser)
	router.GET("users/avatar/:id", middleware.AuthAdminMiddleware(), userWebHandler.FormUplaodAvater)
	router.POST("users/avatar/:id", middleware.AuthAdminMiddleware(), userWebHandler.UploadAvatar)

	//!Router Web Static Campaigns
	router.GET("/campaigns", middleware.AuthAdminMiddleware(), campaignWebHandller.Index)
	router.GET("/campaigns/new", middleware.AuthAdminMiddleware(), campaignWebHandller.FormSelectCreateUser)
	router.POST("/campaigns", middleware.AuthAdminMiddleware(), campaignWebHandller.CreateCampaignUser)
	router.GET("/campaigns/image/:id", middleware.AuthAdminMiddleware(), campaignWebHandller.FormUploadCampaignImage)
	router.POST("/campaigns/image/:id", middleware.AuthAdminMiddleware(), campaignWebHandller.UploadCampaignImage)
	router.GET("/campaigns/edit/:id", middleware.AuthAdminMiddleware(), campaignWebHandller.FormUpdateCampaign)
	router.POST("/campaigns/update/:id", middleware.AuthAdminMiddleware(), campaignWebHandller.UpdateCampaign)
	router.GET("/campaigns/show/:id", middleware.AuthAdminMiddleware(), campaignWebHandller.ShowDetailCampaign)

	//!Router Web Static Transactions
	router.GET("/transactions", middleware.AuthAdminMiddleware(), transactionWebHandler.Index)

	router.Run(sPort) //! default PORT 8080
}

//! input (memasukkan data atau mengirim request dari client) -> Handler (mapping input ke struct) -> memanggil Service (melakukan bisnis proses, mapping struct) -> repository(akses ke database, berupa CRUD) -> memanggil DB
