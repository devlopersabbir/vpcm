package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devlopersabbir/vpcm/internal/inventory"
)

// GetServers returns the list of all registered servers via REST API or direct database connection
func (a *App) GetServers() ([]inventory.ServerView, error) {
	url := fmt.Sprintf("%s/servers", a.getAPIURL())
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err == nil {
		a.addAuthHeader(req)
		resp, err := httpClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var list []inventory.ServerView
				if err := json.NewDecoder(resp.Body).Decode(&list); err == nil {
					return list, nil
				}
			}
		}
	}

	// Fallback to direct database query (SQLite or user's custom MongoDB Connection URI)
	directServers, err := a.getDirectDBServers()
	if err == nil {
		return directServers, nil
	}

	return []inventory.ServerView{}, nil
}

// GetServer returns the detailed view of a single server via REST API
func (a *App) GetServer(id uint) (*inventory.ServerView, error) {
	url := fmt.Sprintf("%s/servers/%d", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	a.addAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var view inventory.ServerView
	if err := json.NewDecoder(resp.Body).Decode(&view); err != nil {
		return nil, err
	}
	return &view, nil
}

// AddServer registers a new server in the inventory via REST API
func (a *App) AddServer(name string, host string, port int, username string, authType string, authSecret string) error {
	url := fmt.Sprintf("%s/servers", a.getAPIURL())
	srv := &inventory.Server{
		Name:       name,
		Host:       host,
		Port:       port,
		Username:   username,
		AuthType:   authType,
		AuthSecret: authSecret,
	}

	data, err := json.Marshal(srv)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	a.addAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}
	return nil
}

// DeleteServer removes a server from the inventory via REST API
func (a *App) DeleteServer(id uint) error {
	url := fmt.Sprintf("%s/servers/%d", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	a.addAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}
	return nil
}
