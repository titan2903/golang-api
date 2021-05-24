package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/middleware"
	"bwastartup/user"
	"log"

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

	//! Users
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	//!Campaigns
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	
	router := gin.Default()
	api := router.Group("/api/v1")

	//! Router access immage
	router.Static("/images", "./images")

	//! Router Users
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailHasBeenRegister)
	api.POST("/avatars", middleware.AuthMiddleware(authService, userService) ,userHandler.UploadAvatar)

	//!Router Campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)

	router.Run()

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