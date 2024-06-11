package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	items   = make(map[int]string)
	nextID  = 1
	itemsMu sync.Mutex
)

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/greet", greetHandler)
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/items/", itemHandler)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "This is my first API using GoLang"}`)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "API is running smoothly"}`)
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"greeting": "Hello, %s!"}`, name)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItems(w, r)
	case http.MethodPost:
		createItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getItem(w, r, id)
	case http.MethodPut:
		updateItem(w, r, id)
	case http.MethodDelete:
		deleteItem(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getItems(w http.ResponseWriter, r *http.Request) {
	itemsMu.Lock()
	defer itemsMu.Unlock()

	itemsList := make([]string, 0, len(items))
	for _, item := range items {
		itemsList = append(itemsList, item)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemsList)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	defer itemsMu.Unlock()

	id := nextID
	nextID++
	items[id] = reqBody.Name

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id": %d, "name": "%s"}`, id, reqBody.Name)
}

func getItem(w http.ResponseWriter, r *http.Request, id int) {
	itemsMu.Lock()
	defer itemsMu.Unlock()

	item, exists := items[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id": %d, "name": "%s"}`, id, item)
}

func updateItem(w http.ResponseWriter, r *http.Request, id int) {
	var reqBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	defer itemsMu.Unlock()

	if _, exists := items[id]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	items[id] = reqBody.Name

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id": %d, "name": "%s"}`, id, reqBody.Name)
}

func deleteItem(w http.ResponseWriter, r *http.Request, id int) {
	itemsMu.Lock()
	defer itemsMu.Unlock()

	if _, exists := items[id]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	delete(items, id)
	w.WriteHeader(http.StatusNoContent)
}
