package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nikhilr/ghabricator/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv, err := server.New()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	fmt.Printf("Ghabricator listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
