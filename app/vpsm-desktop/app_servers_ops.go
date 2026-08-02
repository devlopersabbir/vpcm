package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devlopersabbir/vpcm/internal/inventory"
)

// GetConnectionHistory returns connection history logs for a server via REST API
func (a *App) GetConnectionHistory(id uint) ([]inventory.ConnectionLog, error) {
	url := fmt.Sprintf("%s/servers/%d/history", a.getAPIURL(), id)
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

	if resp.StatusCode == http.StatusNotFound {
		return []inventory.ConnectionLog{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var logs []inventory.ConnectionLog
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// ScanServer triggers background data collection (OS, Hardware, Network, etc.) via REST API
func (a *App) ScanServer(id uint) error {
	url := fmt.Sprintf("%s/servers/%d/scan", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	a.addAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}
	return nil
}

// ToggleFavorite switches the favorite state of a server via REST API
func (a *App) ToggleFavorite(id uint) (bool, error) {
	url := fmt.Sprintf("%s/servers/%d/favorite", a.getAPIURL(), id)
	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return false, err
	}
	a.addAuthHeader(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var res struct {
		IsFavorite bool `json:"is_favorite"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}
	return res.IsFavorite, nil
}
