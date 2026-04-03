package tasks

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var Tasks = []Task{
	{ID: 1, Title: "Task 1", Status: "pending"},
	{ID: 2, Title: "Task 2", Status: "completed"},
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, Tasks)
}
