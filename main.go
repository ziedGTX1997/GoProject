package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Task struct to hold task information
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Global variable to store tasks
var (
	tasks  = []Task{}
	taskID = 1
	mu     sync.Mutex
)

// Handler to add a task
func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mu.Lock()
	task.ID = taskID
	taskID++
	tasks = append(tasks, task)
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Handler to get all tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addTask(w, r)
		} else if r.Method == http.MethodGet {
			getTasks(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}
