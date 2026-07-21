package main

import (
	"context"

	"github.com/devlopersabbir/vpcm/internal/config"
	"github.com/devlopersabbir/vpcm/internal/database"
	"github.com/devlopersabbir/vpcm/internal/inventory"
)

// App struct
type App struct {
	ctx     context.Context
	service inventory.ServerService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = a.initService()
}

func (a *App) initService() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.Database.Driver == "sqlite" {
		db, err := database.InitSQLite(cfg.Database.Path)
		if err != nil {
			return err
		}
		repo, err := inventory.NewSQLiteRepository(db)
		if err != nil {
			return err
		}
		a.service = inventory.NewService(repo)
		return nil
	}

	db, err := database.InitMongo(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		return err
	}
	repo := inventory.NewMongoRepository(db)
	a.service = inventory.NewService(repo)
	return nil
}

// GetServers returns the list of all registered servers
func (a *App) GetServers() ([]inventory.ServerView, error) {
	if a.service == nil {
		if err := a.initService(); err != nil {
			return nil, err
		}
	}
	return a.service.ListServers(a.ctx)
}

// GetServer returns the detailed view of a single server
func (a *App) GetServer(id uint) (*inventory.ServerView, error) {
	if a.service == nil {
		if err := a.initService(); err != nil {
			return nil, err
		}
	}
	return a.service.GetServer(a.ctx, id)
}

// AddServer registers a new server in the inventory
func (a *App) AddServer(name string, host string, port int, username string, authType string, authSecret string) error {
	if a.service == nil {
		if err := a.initService(); err != nil {
			return err
		}
	}
	srv := &inventory.Server{
		Name:       name,
		Host:       host,
		Port:       port,
		Username:   username,
		AuthType:   authType,
		AuthSecret: authSecret,
	}
	return a.service.AddServer(a.ctx, srv)
}

// DeleteServer removes a server from the inventory
func (a *App) DeleteServer(id uint) error {
	if a.service == nil {
		if err := a.initService(); err != nil {
			return err
		}
	}
	return a.service.RemoveServer(a.ctx, id)
}

// GetConnectionHistory returns connection history logs for a server
func (a *App) GetConnectionHistory(id uint) ([]inventory.ConnectionLog, error) {
	if a.service == nil {
		if err := a.initService(); err != nil {
			return nil, err
		}
	}
	return a.service.GetConnectionHistory(a.ctx, id)
}

// ScanServer triggers background data collection (OS, Hardware, Network, etc.)
func (a *App) ScanServer(id uint) error {
	if a.service == nil {
		if err := a.initService(); err != nil {
			return err
		}
	}
	return a.service.ScanInventory(a.ctx, id)
}

// ToggleFavorite switches the favorite state of a server
func (a *App) ToggleFavorite(id uint) (bool, error) {
	if a.service == nil {
		if err := a.initService(); err != nil {
			return false, err
		}
	}
	return a.service.ToggleFavorite(a.ctx, id)
}

// GetConfig returns the active settings
func (a *App) GetConfig() (*config.Config, error) {
	return config.Load()
}

// SaveConfig updates the active settings
func (a *App) SaveConfig(cfg config.Config) error {
	err := config.Save(&cfg)
	if err != nil {
		return err
	}
	return a.initService()
}
