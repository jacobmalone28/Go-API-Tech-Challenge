package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jacob-tech-challenge/api"
	"github.com/jacob-tech-challenge/config"
	"github.com/jacob-tech-challenge/database"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("Startup failed. err: %v", err)
	}
}

func run(ctx context.Context) error {
	// initialize configuration
	cfg, err := config.New()
	if err != nil {
		return err
	}
	// connect to database
	db, err := database.Connect(cfg, sql.Open)

	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database connection. err: %v", err)
		}
	}()

	// initialize router
	r := api.SetupRoutes(db)

	server := &http.Server{
		Addr:    cfg.HTTP_Domain + cfg.HTTP_Port,
		Handler: r,

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server listening on %s", server.Addr)

	// start api server
	log.Println("Starting server on port 8080")
	return server.ListenAndServe()
}
