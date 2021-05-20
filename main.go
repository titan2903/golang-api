package main

import (
	"bwastartup/user"
	"fmt"
	"log"

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

	fmt.Println("Connected...")

	var users []user.User
	length := len(users)
	fmt.Println(length)

	db.Find(&users) //! type harus pointer

	length = len(users)
	fmt.Println(length)

	for _, user := range users {
		fmt.Println(user.Name)
		fmt.Println(user.Email)
	}
}