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

	var invRepo inventory.ServerRepository
	var noteRepo notes.NoteRepository

	if cfg.Database.Driver == "sqlite" {
		db, err := database.InitSQLite(cfg.Database.Path)
		if err != nil {
			log.Fatalf("failed to init SQLite database: %v", err)
		}
		invRepo, err = inventory.NewSQLiteRepository(db)
		if err != nil {
			log.Fatalf("failed to init inventory repo: %v", err)
		}
		noteRepo = notes.NewSQLiteRepository(db)
	} else {
		db, err := database.InitMongo(cfg.Database.URI, cfg.Database.Name)
		if err != nil {
			log.Fatalf("failed to init MongoDB: %v", err)
		}
		invRepo = inventory.NewMongoRepository(db)
		noteRepo = notes.NewMongoRepository(db)
	}

	invSvc := inventory.NewService(invRepo)
	noteSvc := notes.NewService(noteRepo)

	server := api.NewServer(invSvc, noteSvc)

	log.Printf("Starting VPSM API server on %s:%d", cfg.API.Host, cfg.API.Port)
	if err := server.Start(cfg.API.Host, cfg.API.Port); err != nil {
		log.Fatalf("API server failed: %v", err)
	}
}
