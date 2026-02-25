package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nikhilr/ghabricator/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Determine repo root: default to the directory containing this binary's source,
	// or use REPO_ROOT env var.
	repoRoot := os.Getenv("REPO_ROOT")
	if repoRoot == "" {
		// Default: assume running from repo root
		var err error
		repoRoot, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	repoRoot, _ = filepath.Abs(repoRoot)

	srv, err := server.New(repoRoot)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	fmt.Printf("Ghabricator listening on :%s (root: %s)\n", port, repoRoot)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
