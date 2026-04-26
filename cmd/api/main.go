package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yourusername/shorturl/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting API server on port %s", cfg.ServerPort)
	log.Printf("Database: host=%s, port=%d, user=%s, db=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Short URL API is running\n")
	})

	addr := ":" + cfg.ServerPort
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
