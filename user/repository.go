package user

import (
	"gorm.io/gorm"
)


type Repository interface { //! public Interface
	Save(user User) (User, error) //! for save of create user to DB
	FindByEmail(email string) (User, error) //! find user by email in DB 
	FindByID(ID int) (User, error) //! find user by id in DB
	Update(user User) (User, error) //! update entity user in DB
	FindAll() ([]User, error)
}

type repository struct { //! tidak bersifat public
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository { //! membuat object baru dari repository dan nilai db dari repository di isi sesuai parameter di NewRepository
	return &repository{db}
}

func(r *repository) Save(user User) (User, error) { //! create new user
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func(r *repository) FindByEmail(email string) (User, error){ //! find user by email
	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func(r *repository) FindByID(ID int) (User, error) {
	var user User
	err := r.db.Where("id = ?", ID).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func(r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func(r *repository) FindAll() ([]User, error) {
	var users []User

	err := r.db.Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}