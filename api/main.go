package main

import (
	"fmt"
	"github.com/go-audio/audiotools"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No file in request body"))
		return
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error reading file"))
		return
	}

	format, err := audiotools.HeaderFormat(fileContent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error reading file format"))
		return
	}

	if format == "unknown" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid file format"))
		return
	}

	err2 := os.WriteFile(fmt.Sprintf("./%s", fileHeader.Filename), fileContent, 0666)
	if err2 != nil {
		log.Fatal(err2)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Uploaded"))
}
