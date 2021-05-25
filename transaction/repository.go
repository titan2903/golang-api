package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetCampaignByID(campaingID int) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository { //! membuat object baru dari repository dan nilai db dari repository di isi sesuai parameter di NewRepository untuk di panggil di main.go
	return &repository{db}
}

func(r *repository) GetCampaignByID(campaingID int) ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("User").Where("campaign_id = ?", campaingID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	} 

	return transaction, nil
}