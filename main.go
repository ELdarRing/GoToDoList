package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var tasks []Task
var idCounter = 1

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", tasksHandler)

	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		json.NewEncoder(w).Encode(tasks)

	case "POST":
		var task Task
		json.NewDecoder(r.Body).Decode(&task)

		task.ID = idCounter
		idCounter++

		tasks = append(tasks, task)

		json.NewEncoder(w).Encode(task)
	}
}