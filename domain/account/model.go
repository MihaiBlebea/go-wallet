package account

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Account struct {
	ID        string    `json:"id" gorm:"type:string;primaryKey;"`
	MonzoID   string    `json:"monzo_id" gorm:"type:string;"`
	Balance   int       `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(monzoID string, balance int, currency string) *Account {
	return &Account{
		ID:       uuid.NewV4().String(),
		MonzoID:  monzoID,
		Balance:  balance,
		Currency: currency,
	}
}
