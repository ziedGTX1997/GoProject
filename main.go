package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// Handler to update a task
func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = updatedTask.Title
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// Handler to delete a task
func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
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
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			updateTask(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}
