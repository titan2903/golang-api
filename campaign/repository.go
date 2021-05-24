package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaingn, error)
	FindByUserID(userID int) ([]Campaingn, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository { //! membuat object baru dari repository dan nilai db dari repository di isi sesuai parameter di NewRepository
	return &repository{db}
}

func(r *repository) FindAll() ([]Campaingn, error) {
	var campaigns []Campaingn

	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func(r *repository) FindByUserID(userID int) ([]Campaingn, error) {
	var campaigns []Campaingn

	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error //! preload merupakan ngeload relasinya dan mengambil data relasinya
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}