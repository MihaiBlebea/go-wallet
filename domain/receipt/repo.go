package receipt

import (
	"errors"

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
	Save(receipt *Receipt) (string, error)
	Delete(id string) error
	All() ([]Receipt, error)
	FindByShortID(shortID string) (*Receipt, error)
	FindByTransactionID(transactionID string) ([]Receipt, error)
}

func NewRepo(conn *gorm.DB) Repo {
	return &repo{conn}
}

func (r *repo) Save(receipt *Receipt) (string, error) {
	cmd := r.conn.Create(receipt)
	if cmd.RowsAffected == 0 {
		return "", ErrNotInserted
	}

	return receipt.ID, cmd.Error
}

func (r *repo) Delete(id string) error {
	receipt := Receipt{}
	if err := r.conn.Where("id = ?", id).Find(&receipt).Error; err != nil {
		return err
	}

	cmd := r.conn.Delete(receipt)
	if cmd.RowsAffected == 0 {
		return ErrNoRecord
	}

	return cmd.Error
}

func (r *repo) FindByShortID(shortID string) (*Receipt, error) {
	receipt := Receipt{}
	err := r.conn.Where("short_id = ?", shortID).Find(&receipt).Error
	if err != nil {
		return &Receipt{}, err
	}

	return &receipt, nil
}

func (r *repo) FindByTransactionID(transactionID string) ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	err := r.conn.Where("transaction_id = ?", transactionID).Find(&receipts).Error
	if err != nil {
		return []Receipt{}, err
	}

	return receipts, nil
}

func (r *repo) All() ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	err := r.conn.Find(&receipts).Error
	if err != nil {
		return []Receipt{}, err
	}

	return receipts, nil
}
