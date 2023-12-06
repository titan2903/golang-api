package healthcheck

import "time"

type Service interface {
	HealthcheckService() (Healthcheck, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) HealthcheckService() (Healthcheck, error) {
	// You need to implement the logic to perform the health check here
	// For demonstration purposes, I'm creating a simple Healthcheck with static values.
	check := Healthcheck{
		ServiceName: "Golang Crowdfunding Services",
		Status:      "OK",
		Description: "Everything is running well",
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}
	return check, nil
}
