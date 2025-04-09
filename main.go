package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting server...")

	serverMux := http.NewServeMux()
	server := &http.Server{
		Addr: ":8080",
	}

	err := http.ListenAndServe(server.Addr, serverMux)
	if err != nil {
		fmt.Println("Error starting server", err)
		return
	} else {
		fmt.Println("Server started successfully")
		fmt.Println("Listening on port 8080")
	}
}
