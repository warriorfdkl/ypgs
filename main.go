package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"task"`
	Done bool   `json:"done"`
}

var (
	tasks  []Task
	mu     sync.Mutex
	nextID = 1
)

func main() {
	http.HandleFunc("/set-task", createHandler)
	http.HandleFunc("/tasks", listHandler)
	http.HandleFunc("/task/", patchHandler)
	http.HandleFunc("/tasks/", deleteHandler)
	fmt.Println("Сервер на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST only", 405)
		return
	}
	var body struct {
		Task string `json:"task"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	mu.Lock()
	id := nextID
	nextID++
	tasks = append(tasks, Task{ID: id, Text: body.Task, Done: false})
	mu.Unlock()

	w.WriteHeader(201)
	fmt.Fprintf(w, "%d", id)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET only", 405)
		return
	}
	mu.Lock()
	json.NewEncoder(w).Encode(tasks)
	mu.Unlock()
}

func patchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "PATCH only", 405)
		return
	}
	idStr := r.URL.Path[len("/task/"):]
	id, _ := strconv.Atoi(idStr)
	var upd struct {
		Task *string `json:"task,omitempty"`
		Done *bool   `json:"done,omitempty"`
	}
	json.NewDecoder(r.Body).Decode(&upd)

	mu.Lock()
	defer mu.Unlock()
	for i := range tasks {
		if tasks[i].ID == id {
			if upd.Task != nil {
				tasks[i].Text = *upd.Task
			}
			if upd.Done != nil {
				tasks[i].Done = *upd.Done
			}
			w.WriteHeader(200)
			return
		}
	}
	http.Error(w, "not found", 404)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "DELETE only", 405)
		return
	}
	idStr := r.URL.Path[len("/tasks/"):]
	id, _ := strconv.Atoi(idStr)

	mu.Lock()
	defer mu.Unlock()
	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(200)
			return
		}
	}
	http.Error(w, "not found", 404)
}
