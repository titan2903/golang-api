package campaign

type Service interface {
	GetCampaigns(UserID int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	//! melakukan switching percabangan
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	} else {
		campaigns, err := s.repository.FindAll()
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
}