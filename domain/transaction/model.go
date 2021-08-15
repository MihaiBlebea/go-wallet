package transaction

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Transaction struct {
	ID          string    `json:"id" gorm:"type:string;primaryKey"`
	MonzoID     string    `json:"monzo_id"`
	MerchantID  string    `json:"merchant_id"`
	Completed   bool      `json:"completed"`
	Amount      int       `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Notes       string    `json:"notes"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func New(
	monzoID,
	merchantID string,
	amount int,
	currency string,
	description string,
	notes string,
	category string) *Transaction {

	return &Transaction{
		ID:          uuid.NewV4().String(),
		MonzoID:     monzoID,
		MerchantID:  merchantID,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Notes:       notes,
		Category:    category,
	}
}
