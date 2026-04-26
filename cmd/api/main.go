package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/marifyahya/shorturl-generator-app/internal/config"
	"github.com/marifyahya/shorturl-generator-app/internal/handler"
	"github.com/marifyahya/shorturl-generator-app/internal/repository"
	"github.com/marifyahya/shorturl-generator-app/internal/service"
)

func main() {
	// 0. Parse Flags for Migration Control
	migrateFlag := flag.String("migrate", "up", "Direction of migration: up or down")
	flag.Parse()

	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Database
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Run Migrations based on flag
	if *migrateFlag == "down" {
		if err := repository.RollbackMigration(db, cfg.DBName); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		log.Println("Successfully rolled back 1 step. Exiting...")
		return
	}

	if err := repository.RunMigrations(db, cfg.DBName); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Initialize Repository, Service, and Handler
	repo := repository.NewPostgresURLRepository(db)
	svc := service.NewURLService(repo)
	h := handler.NewURLHandler(svc, cfg)

	// 4. Setup Routing
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("/api/shorten", h.Shorten)
	mux.HandleFunc("/api/stats/", h.GetStats)

	// Redirection (Catch-all for short codes)
	mux.HandleFunc("/", h.Redirect)

	// 5. Start Server
	log.Printf("Server started on port %s", cfg.ServerPort)
	log.Printf("Base URL: %s", cfg.BaseURL)

	addr := ":" + cfg.ServerPort
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
