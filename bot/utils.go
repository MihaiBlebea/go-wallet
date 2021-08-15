package bot

import (
	"errors"
	"strconv"
	"strings"
)

type ReceiptPayload struct {
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}

func (s *service) parseReceiptCommand(message string) (*ReceiptPayload, error) {
	payload := ReceiptPayload{}
	if strings.Contains(message, OnAdd) == false {
		return &payload, errors.New("The message does not contain the command expected")
	}

	data := strings.Split(strings.TrimSpace(strings.ReplaceAll(message, OnAdd, "")), " ")
	if len(data) > 4 {
		return &payload, errors.New("Too many arguments")
	}

	amount, err := s.toIntAmount(data[1])
	if err != nil {
		return &payload, err
	}

	quantity, err := strconv.Atoi(data[2])
	if err != nil {
		return &payload, err
	}

	payload.Description = data[0]
	payload.Amount = amount
	payload.Currency = "GBP"
	payload.Quantity = quantity
	if len(data) == 4 {
		payload.Unit = data[3]
	} else {
		payload.Unit = "unit"
	}

	return &payload, nil
}

func (s *service) parseDeleteCommand(message string) (string, error) {
	if strings.Contains(message, OnDelete) == false {
		return "", errors.New("The message does not contain the command expected")
	}

	data := strings.Split(strings.TrimSpace(strings.ReplaceAll(message, OnDelete, "")), " ")
	if len(data) > 1 {
		return "", errors.New("Too many arguments")
	}

	return data[0], nil
}

func (s *service) toIntAmount(amount string) (int, error) {

	// Check it has . if not add 00
	if !strings.Contains(amount, ".") {
		amount += ".00"
	}

	// Check it has two digits after . if not add 0
	if strings.Index(amount, ".") > len(amount)-3 {
		amount += "0"
	}

	// Remove dot to get cents
	amount = strings.Replace(amount, ".", "", 1)

	// Remove any stray formatting users might add
	amount = strings.Trim(amount, "£€$ ")

	// Parse int
	return strconv.Atoi(amount)
}
