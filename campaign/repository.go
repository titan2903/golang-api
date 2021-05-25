package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	SaveCampaign(campaign Campaign) (Campaign, error)
	UpdateCampaign(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository { //! membuat object baru dari repository dan nilai db dari repository di isi sesuai parameter di NewRepository
	return &repository{db}
}

func(r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func(r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Find(&campaigns).Error //! preload merupakan ngeload relasinya dan mengambil data relasinya
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func(r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign

	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error //! melakukan preload untuk mengambil data table CampaignImage dan User
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func(r *repository) SaveCampaign(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func(r *repository) UpdateCampaign(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}