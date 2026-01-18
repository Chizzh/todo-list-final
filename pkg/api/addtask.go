package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"todo-list-final/pkg/db"
)

type jsonError struct {
	Error string `json:"error"`
}

type jsonID struct {
	ID string `json:"id"`
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, jsonError{err.Error()})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, jsonError{"title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, jsonError{err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, jsonError{err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, jsonID{strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты")
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("DateFormat")
		} else {
			task.Date = next
		}
	}

	return nil
}
