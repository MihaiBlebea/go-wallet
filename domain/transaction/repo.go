package transaction

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNoRecord    error = errors.New("Record not found")
	ErrNoRecords   error = errors.New("Records not found")
	ErrNotInserted error = errors.New("Could not insert record")
)

type repo struct {
	conn *gorm.DB
}

type Repo interface {
	Save(transaction *Transaction) (string, error)
	MarkCompleted(id string) error
	FindFirstUncompleted() (*Transaction, error)
	FindInInterval(start, end time.Time) ([]Transaction, error)
}

func NewRepo(conn *gorm.DB) Repo {
	return &repo{conn}
}

func (r *repo) Save(transaction *Transaction) (string, error) {
	cmd := r.conn.Create(transaction)
	if cmd.RowsAffected == 0 {
		return "", ErrNotInserted
	}

	return transaction.ID, cmd.Error
}

func (r *repo) MarkCompleted(id string) error {
	cmd := r.conn.Model(&Transaction{}).Where("id = ?", id).Update("completed", true)
	if cmd.RowsAffected == 0 {
		return ErrNoRecord
	}

	return cmd.Error
}

func (r *repo) FindFirstUncompleted() (*Transaction, error) {
	transaction := Transaction{}
	err := r.conn.Where("completed = ?", false).First(&transaction).Error
	if err != nil {
		return &transaction, err
	}

	return &transaction, nil
}

func (r *repo) FindInInterval(start, end time.Time) ([]Transaction, error) {
	transactions := make([]Transaction, 0)
	err := r.conn.Where("created_at > ? AND created_at < ?", start, end).Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
