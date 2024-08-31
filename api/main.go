package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", handleUpload)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	err := r.ParseMultipartForm(1064)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error parsing request body"))
		return
	}

	fileHeaders, fileFieldIncluded := r.MultipartForm.File["file"]
	if !fileFieldIncluded {
		w.WriteHeader(400)
		w.Write([]byte("Missing file field in form data"))
		return
	}

	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error opening file"))
			return
		}
		defer file.Close()

		fileContent, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error reading file"))
			return
		}

		err2 := os.WriteFile(fmt.Sprintf("./%s", fileHeader.Filename), fileContent, 0666)
		if err2 != nil {
			log.Fatal(err)
		}

	}

	w.WriteHeader(201)
	w.Write([]byte("Uploaded"))
}
