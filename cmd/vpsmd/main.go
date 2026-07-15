package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
	"github.com/devlopersabbir/vpcm/internal/logger"
	"github.com/devlopersabbir/vpcm/internal/scheduler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	logger.Init(cfg.Logging.Level, cfg.Logging.Format, os.Stdout)

	_, err = database.Init(cfg.Database.Path)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	sched := scheduler.New()

	// Skeleton job to check in every 30 seconds
	_, err = sched.AddJob("*/30 * * * * *", func() {
		log.Println("VPSM Daemon background check tick...")
	})
	if err != nil {
		log.Fatalf("failed to schedule job: %v", err)
	}

	sched.Start()

	// Graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	sched.Stop()
	log.Println("VPSM Daemon stopped gracefully")
}
