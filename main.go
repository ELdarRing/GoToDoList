package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var tasks []Task
var idCounter = 1
const fileName = "tasks.json"

func main() {
	loadTasks()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/toggle", toggleHandler)

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
		saveTasks()

		json.NewEncoder(w).Encode(task)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	saveTasks()
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = !tasks[i].Done
			break
		}
	}

	saveTasks()
}

func saveTasks() {
	data, _ := json.Marshal(tasks)
	ioutil.WriteFile(fileName, data, 0644)
}

func loadTasks() {
	if _, err := os.Stat(fileName); err == nil {
		data, _ := ioutil.ReadFile(fileName)
		json.Unmarshal(data, &tasks)

		for _, t := range tasks {
			if t.ID >= idCounter {
				idCounter = t.ID + 1
			}
		}
	}
}