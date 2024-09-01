package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nicomellon/auto-master/services"
)

type ResponseBody struct {
	Data   *Resource `json:"data,omitempty"`
	Errors []Error   `json:"errors,omitempty"`
}

type Resource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Error struct {
	Detail string `json:"detail"`
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	if r.Method != "POST" {
		payload := ResponseBody{nil, []Error{Error{"Method not allowed"}}}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(payload)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		payload := ResponseBody{nil, []Error{Error{"No file in request body"}}}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload)
		return
	}

	uploadService := services.NewUploadService()
	var payload ResponseBody
	var status int
	if id, err := uploadService.Upload(file); err != nil {
		payload = ResponseBody{nil, []Error{Error{fmt.Sprint(err)}}}
		status = http.StatusBadRequest
	} else {
		payload = ResponseBody{&Resource{id, "uploads"}, nil}
		status = http.StatusCreated
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
