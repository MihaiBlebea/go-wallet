package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	b, _ := json.Marshal(resp)

	w.Write(b)
}

func stringToTime(value string) (time.Time, error) {
	if strings.Contains(value, "-") == false {
		return time.Time{}, errors.New("Invalid date format. Please use YYYY-MM-DD")
	}

	data := strings.Split(value, "-")
	if len(data) != 3 {
		return time.Time{}, errors.New("Invalid date format. Please use YYYY-MM-DD")
	}

	if len(data[0]) != 4 && len(data[1]) != 2 && len(data[2]) != 2 {
		return time.Time{}, errors.New("Invalid date format. Please use YYYY-MM-DD")
	}

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return time.Time{}, err
	}

	yearInt, err := strconv.Atoi(data[0])
	if err != nil {
		return time.Time{}, err
	}

	monthInt, err := strconv.Atoi(data[1])
	if err != nil {
		return time.Time{}, err
	}

	dayInt, err := strconv.Atoi(data[2])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, loc), nil
}
