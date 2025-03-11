package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var tasks []Task

// Load tasks from task.json
func loadTasks() {
	file, err := ioutil.ReadFile("task.json")
	if err != nil {
		log.Println("Could not read task.json:", err)
		return
	}
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		log.Println("Could not unmarshal task.json:", err)
	}
}

// Save tasks to task.json
func saveTasks() {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Println("Could not marshal tasks:", err)
		return
	}
	err = ioutil.WriteFile("task.json", data, 0644)
	if err != nil {
		log.Println("Could not write to task.json:", err)
	}
}

// Get all tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	loadTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Get task by ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	loadTasks()
	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.NotFound(w, r)
}

// Create a new task
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	loadTasks()
	newTask.ID = len(tasks) + 1 // Simple auto-increment ID
	tasks = append(tasks, newTask)
	saveTasks()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

// Update a task by ID
func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	loadTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = updatedTask.Title
			saveTasks()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.NotFound(w, r)
}

// Delete a task by ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	loadTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

// Simulate a long-running process
func processTask(w http.ResponseWriter, r *http.Request) {
	// Send immediate response to the client
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Task processing started in the background.\n"))

	// Run the long-running task in a goroutine
	go func() {
		fmt.Println("Processing task... This will take 5 seconds.")
		time.Sleep(5 * time.Second)
		fmt.Println("Task processing completed.")
	}()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", getTaskByID).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	r.HandleFunc("/taskse/process", processTask).Methods("GET") // New route

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
