package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MihaiBlebea/go-wallet/bot"
	"github.com/MihaiBlebea/go-wallet/domain"
	"github.com/MihaiBlebea/go-wallet/domain/transaction"
)

type TransactionCreatedRequest struct {
	Type string `json:"type"`
	Data struct {
		AccountID   string `json:"account_id"`
		Amount      int    `json:"amount"`
		Created     string `json:"created"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Category    string `json:"category"`
		IsLoad      bool   `json:"is_load"`
		Settled     string `json:"settled"`
		Notes       string `json:"notes"`
		Merchant    struct {
			ID      string `json:"id"`
			Address struct {
				Address   string  `json:"address"`
				City      string  `json:"city"`
				Country   string  `json:"country"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				Postcode  string  `json:"postcode"`
				Region    string  `json:"region"`
			} `json:"address"`
			Logo     string `json:"logo"`
			Emoji    string `json:"emoji"`
			Name     string `json:"name"`
			Category string `json:"category"`
		} `json:"merchant"`
	} `json:"data"`
}

type TransactionCreatedResponse struct {
	ID      string `json:"id,omitempty"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func WebhookEndpoint(wallet domain.Wallet, bot bot.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := TransactionCreatedRequest{}
		response := TransactionCreatedResponse{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		if req.Data.IsLoad == true {
			return
		}

		transaction, err := wallet.StoreTransaction(
			req.Data.ID,
			req.Data.Merchant.ID,
			req.Data.Amount,
			req.Data.Currency,
			req.Data.Description,
			req.Data.Notes,
			req.Data.Category,
		)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		askUserForReceipt(bot, transaction)

		response.Success = true
		response.ID = transaction.ID

		sendResponse(w, response, 200)
	})
}

func askUserForReceipt(bot bot.Service, t *transaction.Transaction) {
	message := fmt.Sprintf(
		"Can you update this transaction of %s from %s, please?",
		bot.ToAmount(t.Amount, t.Currency),
		t.Description,
	)
	bot.SendMessageToUser(message)
	bot.SendMessageToUser("Format /add <product> <fee> <quantity> <unit?>")
}
