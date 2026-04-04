package main

import (
	"fmt"
	"net/http"
	"task-api/tasks"
)

func getHelloWorld(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		App    string `json:"app"`
		Author string `json:"author"`
		Repo   string `json:"repo"`
	}

	body := Body{
		App:    "Task API Go",
		Author: "Tegar Wijaya Kusuma",
		Repo:   "https://github.com/tgr-wjya/go-task-api",
	}

	tasks.WriteJSON(w, http.StatusOK, body)
}

func getPlain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Hello, World")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", getHelloWorld)
	mux.HandleFunc("GET /plain", getPlain)
	mux.HandleFunc("GET /tasks/all", tasks.GetAll)

	fmt.Println("Server listening at 8080")
	http.ListenAndServe(":8080", mux)
}
