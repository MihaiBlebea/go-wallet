package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/MihaiBlebea/go-wallet/domain"
	"github.com/MihaiBlebea/go-wallet/domain/report"
)

type ReportRequest struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type ReportResponse struct {
	Date    *report.Report `json:"data,omitempty"`
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
}

func ReportEndpoint(wallet domain.Wallet) http.Handler {
	parse := func(r *http.Request) (*ReportRequest, error) {
		req := ReportRequest{}

		start := r.FormValue("start")
		end := r.FormValue("end")

		if start == "" || end == "" {
			return &req, errors.New("Invalid params")
		}

		startTime, err := stringToTime(start)
		if err != nil {
			return &req, err
		}
		endTime, err := stringToTime(end)
		if err != nil {
			return &req, err
		}

		req.Start = startTime
		req.End = endTime

		return &req, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ReportResponse{}

		request, err := parse(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		report, err := wallet.Report(request.Start, request.End)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Date = report
		response.Success = true

		sendResponse(w, response, 200)
	})
}
