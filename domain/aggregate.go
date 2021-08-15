package domain

import (
	"errors"
	"time"

	"github.com/MihaiBlebea/go-wallet/domain/account"
	"github.com/MihaiBlebea/go-wallet/domain/receipt"
	"github.com/MihaiBlebea/go-wallet/domain/report"
	"github.com/MihaiBlebea/go-wallet/domain/transaction"
	"gorm.io/gorm"
)

type Wallet interface {
	StoreTransaction(
		monzoID,
		merchantID string,
		amount int,
		currency string,
		description string,
		notes string,
		category string) (*transaction.Transaction, error)
	StoreReceipt(
		description string,
		amount int,
		currency string,
		quantity int,
		unit string) (string, int, error)
	DeleteReceipt(shortID string) error
	GetUncompletedTransaction() (*transaction.Transaction, error)
	MarkTransactionCompleted() error
	Report(start, end time.Time) (*report.Report, error)
}

type wallet struct {
	accountRepo     account.Repo
	receiptRepo     receipt.Repo
	transactionRepo transaction.Repo
}

func New(conn *gorm.DB) Wallet {
	return &wallet{
		accountRepo:     account.NewRepo(conn),
		receiptRepo:     receipt.NewRepo(conn),
		transactionRepo: transaction.NewRepo(conn),
	}
}

func (w *wallet) StoreTransaction(
	monzoID,
	merchantID string,
	amount int,
	currency string,
	description string,
	notes string,
	category string) (*transaction.Transaction, error) {

	if amount < 0 {
		amount = -amount
	}

	t := transaction.New(
		monzoID, merchantID, amount, currency, description, notes, category,
	)
	_, err := w.transactionRepo.Save(t)
	if err != nil {
		return &transaction.Transaction{}, err
	}

	return t, nil
}

func (w *wallet) StoreReceipt(
	description string,
	amount int,
	currency string,
	quantity int,
	unit string) (string, int, error) {

	transaction, err := w.transactionRepo.FindFirstUncompleted()
	if err != nil {
		return "", 0, err
	}

	// Validate receipt amount and transaction amount
	receipts, err := w.receiptRepo.FindByTransactionID(transaction.ID)
	if err != nil {
		return "", 0, err
	}

	total := transaction.Amount
	for _, r := range receipts {
		total -= r.Amount
	}

	if total-amount < 0 {
		return "", 0, errors.New("Receipt amount exceeds the transaction amount")
	}

	total -= amount

	rec := receipt.New(
		transaction.ID,
		description,
		amount,
		currency,
		quantity,
		unit,
	)
	_, err = w.receiptRepo.Save(rec)
	if err != nil {
		return "", 0, err
	}

	if total == 0 {
		err := w.transactionRepo.MarkCompleted(transaction.ID)
		if err != nil {
			return "", 0, err
		}
	}

	return rec.ShortID, total, nil
}

func (w *wallet) DeleteReceipt(shortID string) error {
	receipt, err := w.receiptRepo.FindByShortID(shortID)
	if err != nil {
		return err
	}

	return w.receiptRepo.Delete(receipt.ID)
}

func (w *wallet) GetUncompletedTransaction() (*transaction.Transaction, error) {
	return w.transactionRepo.FindFirstUncompleted()
}

func (w *wallet) MarkTransactionCompleted() error {
	transaction, err := w.transactionRepo.FindFirstUncompleted()
	if err != nil {
		return err
	}

	return w.transactionRepo.MarkCompleted(transaction.ID)
}

func (w *wallet) Report(start, end time.Time) (*report.Report, error) {
	transactions, err := w.transactionRepo.FindInInterval(start, end)
	if err != nil {
		return &report.Report{}, err
	}

	report := report.NewWithTransactions(transactions)
	for _, t := range transactions {
		receipts, err := w.receiptRepo.FindByTransactionID(t.ID)
		if err != nil {
			continue
		}

		report.AddReceipts(receipts)
	}

	return report, nil
}
