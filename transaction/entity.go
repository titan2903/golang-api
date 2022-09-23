package transaction

import (
	"golang-api-crowdfunding/campaign"
	"golang-api-crowdfunding/user"
	"time"

	"github.com/leekchan/accounting"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       user.User
	Campaign   campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (t Transaction) AmountFormatIDR() string {
	tr := accounting.Accounting{
		Symbol:    "Rp",
		Precision: 2,
		Thousand:  ".",
		Decimal:   ",",
	}
	return tr.FormatMoney(t.Amount)
}
