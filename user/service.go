package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)


type Service interface { //! bisnis logic
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(input FormUpdateUserInput) (User, error)
}

type service struct { //! memanggil repository
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)

	if err != nil {
		return newUser ,err
	}

	return newUser, nil
}

// mapping struct input ke struct user
// simpan struct User melalui repository

func(s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return user, err
	}

	return user, nil
}

func(s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email

	user, _ := s.repository.FindByEmail(email)

	if user.ID == 0 { //! email tidak di temukan atau bisa di daftarkan
		return true, nil
	}

	return false, nil //! email sudah di gunakan
}

func(s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	/*
		dapatkan user by id
		update attribute avatar file name
		simpan perubahan avatar file name ke DB
	*/

	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation //! belum di simpan ke DB

	updatedUser, err := s.repository.Update(user) //! simpan hasil update ke DB
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func(s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with that ID")
	}

	return user, nil
}

func(s *service) GetAllUsers() ([]User, error) {
	users, err := s.repository.FindAll()
	if err != nil {
		return users, err
	}

	return users, nil
}

func(s *service) UpdateUser(input FormUpdateUserInput) (User, error) {
	user, err := s.repository.FindByID(input.ID)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}