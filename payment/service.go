package payment

import (
	"bwastartup/user"
	"log"
	"strconv"

	"github.com/joho/godotenv"
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

func(s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverKey := myEnv["SERVER_KEY"]
	clientKey := myEnv["CLIENT_KEY"]
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
			OrderID: strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenRes, err := snapGateway.GetToken(snapReq)
	if err!= nil {
		return "", err
	}

	return snapTokenRes.RedirectURL, err
}


/*
	error import cyclec not allowed
		- saling mengimport antar package tidak di bolehkan di golang
*/