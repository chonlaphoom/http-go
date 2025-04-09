package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("starting server...")

	port := "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
