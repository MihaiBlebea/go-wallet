package account

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
	Save(record *Account) (string, error)
	First() (*Account, error)
}

func NewRepo(conn *gorm.DB) Repo {
	return &repo{conn}
}

func (r *repo) First() (*Account, error) {
	account := Account{}
	err := r.conn.First(&account).Error
	if err != nil {
		return &account, err
	}

	if account.ID == "" {
		return &account, ErrNoRecord
	}

	return &account, nil
}

func (r *repo) Save(account *Account) (string, error) {
	cmd := r.conn.Create(account)
	if cmd.RowsAffected == 0 {
		return "", ErrNotInserted
	}

	return account.ID, cmd.Error
}
