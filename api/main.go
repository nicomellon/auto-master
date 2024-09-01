package main

import (
	"log"
	"net/http"

	"github.com/nicomellon/auto-master/handlers"
)

func main() {
	http.HandleFunc("/upload", handlers.UploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
