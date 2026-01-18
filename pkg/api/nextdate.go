package api

import (
	"time"
	"net/http"
	"strings"
	"fmt"
	"strconv"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	now  = time.Date(now.Year(),  now.Month(),  now.Day(),  0, 0, 0, 0, time.UTC)
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat")
	}
	
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid format")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

		case "d":
		if len(parts) != 2 {
        	return "", fmt.Errorf("invalid format")
    	}
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid day interval")
		}
		if n < 1 {
			return "", fmt.Errorf("day interval must be >= 1")
		}
		if n > 400 {
			return "", fmt.Errorf("day interval must be <= 400")
		}

		for {
			date = date.AddDate(0, 0, n)
			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat")
	}

	return date.Format(DateFormat), nil
	
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	res, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(res)); err != nil {
		return
	}
}