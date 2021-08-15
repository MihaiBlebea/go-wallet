package report

import (
	"github.com/MihaiBlebea/go-wallet/domain/receipt"
	"github.com/MihaiBlebea/go-wallet/domain/transaction"
)

type Report struct {
	Transactions []transaction.Transaction
	Receipts     []receipt.Receipt
}

func New(transactions []transaction.Transaction, receipts []receipt.Receipt) *Report {
	return &Report{Transactions: transactions, Receipts: receipts}
}

func NewWithTransactions(transactions []transaction.Transaction) *Report {
	return &Report{Transactions: transactions}
}

func (r *Report) AddReceipts(receipts []receipt.Receipt) {
	r.Receipts = append(r.Receipts, receipts...)
}
