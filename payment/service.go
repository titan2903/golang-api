package payment

import (
	"bwastartup/helper"
	"bwastartup/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {

	serverKey := helper.GoDotEnvVariable("SERVER_KEY")
	clientKey := helper.GoDotEnvVariable("CLIENT_KEY")
	midclient := midtrans.NewClient()
	midclient.ServerKey = serverKey //! my server key
	midclient.ClientKey = clientKey //! my client key
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenRes, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenRes.RedirectURL, err
}

/*
	error import cyclec not allowed
		- saling mengimport antar package tidak di bolehkan di golang
*/
