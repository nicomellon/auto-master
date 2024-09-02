package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicomellon/auto-master/services"
)

type Document struct {
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

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	file, _, err := r.FormFile("file")
	if err != nil {
		payload := Document{nil, []Error{Error{"No file in request body"}}}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload)
		return
	}

	uploadService := services.NewUploadService()
	track, err := uploadService.Upload(file)
	if err != nil {
		payload := Document{nil, []Error{Error{fmt.Sprint(err)}}}
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload)
		return
	}

	payload := Document{&Resource{track.ID.String(), "uploads"}, nil}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payload)
}
