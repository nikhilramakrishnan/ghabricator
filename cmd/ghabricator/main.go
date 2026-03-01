package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nikhilr/ghabricator/internal/server"
)

func main() {
	loadEnvFile(".env")

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

// loadEnvFile reads KEY=VALUE lines from a file and sets them as env vars.
// Silently skips if the file doesn't exist. Does not override existing env vars.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		// Strip surrounding quotes.
		if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
			v = v[1 : len(v)-1]
		}
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}
