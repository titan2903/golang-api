package user

import (
	"gorm.io/gorm"
)


type Repository interface { //! bersifat public
	Save(user User) (User, error) //! untuk save user ke DB
}

type repository struct { //! tidak bersifat public
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository { //! membuat object baru dari repository dan nilai db dari repository di isi sesuai parameter di NewRepository
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}