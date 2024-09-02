package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nicomellon/auto-master/handler"
	"github.com/nicomellon/auto-master/middleware"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /upload", handler.UploadHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Server listening on port", server.Addr)
	log.Fatal(server.ListenAndServe())
}
