package transaction

import (
	"bwastartup/campaign"
	"errors"
)


type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {
	/*
		get campaign
		check campaig.user_id != user_id yang melakukan request
	*/

	//! melakukan auth user
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("not the owner of the campaign's transactions")
	}

	transaction, err := s.repository.GetCampaignByID(input.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}