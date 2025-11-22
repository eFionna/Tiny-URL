package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eFionna/Tiny-URL/internal/app"
	"github.com/eFionna/Tiny-URL/internal/config"
)

func main() {
	cfg := config.Load()

	appInstance, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", appInstance.HandleIndex)
	mux.HandleFunc("/s/", appInstance.HandleRedirect)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server running on port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
