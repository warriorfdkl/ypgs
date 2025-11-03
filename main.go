package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var task string

type RequestBody struct {
	Task string `json:"task"`
}

func main() {
	http.HandleFunc("/set-task", postHandler)
	http.HandleFunc("/get-task", getHandler)
	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST only", 405)
		return
	}
	var b RequestBody
	json.NewDecoder(r.Body).Decode(&b)
	task = b.Task
	fmt.Fprint(w, task)
}
func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET only", 405)
		return
	}
	if task == "" {
		fmt.Fprint(w, "hello,")
	} else {
		fmt.Fprintf(w, "hello, %s", task)
	}
}
