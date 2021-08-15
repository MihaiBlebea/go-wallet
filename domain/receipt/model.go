package receipt

import (
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Receipt struct {
	ID            string    `json:"id" gorm:"type:string;primaryKey"`
	ShortID       string    `json:"short_id"`
	TransactionID string    `json:"transaction_id"`
	Description   string    `json:"description"`
	Amount        int       `json:"amount"`
	Currency      string    `json:"currency"`
	Quantity      int       `json:"quantity"`
	Unit          string    `json:"unit"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func New(transactionID, description string, amount int, currency string, quantity int, unit string) *Receipt {
	return &Receipt{
		ID:            uuid.NewV4().String(),
		ShortID:       generateShortID(),
		TransactionID: transactionID,
		Description:   description,
		Amount:        amount,
		Currency:      currency,
		Quantity:      quantity,
		Unit:          unit,
	}
}

func generateShortID() string {
	rand.Seed(time.Now().UnixNano())
	var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 6)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return fmt.Sprintf("PAY-%s", string(b))
}
