package handler

import (
	"net/http"
)

func HealthEndpoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response struct {
			OK bool `json:"ok"`
		}
		response.OK = true

		sendResponse(w, &response, http.StatusOK)
	})
}
