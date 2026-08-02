package main

import (
	"context"
	"net/http"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
	"github.com/devlopersabbir/vpcm/internal/inventory"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     5 * time.Second,
	},
}

// App struct
type App struct {
	ctx              context.Context
	isTerminalOnly   bool
	terminalServerID uint
	terminalParams   SSHConnectionParams
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) getAPIURL() string {
	cfg, err := config.Load()
	if err == nil && cfg.API.GlobalURL != "" {
		return cfg.API.GlobalURL
	}
	return "http://127.0.0.1:8080"
}

func (a *App) addAuthHeader(req *http.Request) {
	if req == nil {
		return
	}
	cfg, err := config.Load()
	if err == nil && cfg != nil {
		if cfg.API.GuardPassword != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.API.GuardPassword)
			req.Header.Set("X-Cloud-Auth", cfg.API.GuardPassword)
		} else if cfg.API.Token != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.API.Token)
			req.Header.Set("X-Cloud-Auth", cfg.API.Token)
		}
	}
}

func (a *App) getLocalSQLiteServers() ([]inventory.ServerView, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	db, err := database.InitSQLite(cfg.Database.Path)
	if err != nil {
		return nil, err
	}
	repo, err := inventory.NewSQLiteRepository(db)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return repo.ListServerViews(ctx)
}
