package main

import (
	"log"
	"os"

	"github.com/devlopersabbir/vpcm/internal/api"
	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
	"github.com/devlopersabbir/vpcm/internal/inventory"
	"github.com/devlopersabbir/vpcm/internal/logger"
	"github.com/devlopersabbir/vpcm/internal/notes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	logger.Init(cfg.Logging.Level, cfg.Logging.Format, os.Stdout)

	db, err := database.Init(cfg.Database.Path)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// Auto migrate schema
	if err := db.AutoMigrate(&inventory.Server{}, &inventory.Tag{}, &inventory.Software{}, &notes.Note{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	invRepo := inventory.NewRepository(db)
	invSvc := inventory.NewService(invRepo)

	noteRepo := notes.NewRepository(db)
	noteSvc := notes.NewService(noteRepo)

	server := api.NewServer(invSvc, noteSvc)

	log.Printf("Starting VPSM API server on %s:%d", cfg.API.Host, cfg.API.Port)
	if err := server.Start(cfg.API.Host, cfg.API.Port); err != nil {
		log.Fatalf("API server failed: %v", err)
	}
}
