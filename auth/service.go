package auth

import "github.com/dgrijalva/jwt-go"

type Service interface {
	GenerateToken(userID int) (string, error) //! data yang ingin di generate
}

type jwtService struct {
}

var SECRET_KEY = []byte("BANANA_1234")

func NewService() *jwtService { //! bisa memanggil generate token dari package mana pun
	return &jwtService{}
}

func(s *jwtService) GenerateToken(userID int) (string, error) {
	payload := jwt.MapClaims{}
	payload["user_id"] = userID //! data yang ingin di masukkan ke token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload) //! generate token

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}