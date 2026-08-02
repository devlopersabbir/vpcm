package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/devlopersabbir/vpcm/internal/config"
)

// GetConfig returns the active settings
func (a *App) GetConfig() (*config.Config, error) {
	return config.Load()
}

// SaveConfig updates the active settings
func (a *App) SaveConfig(cfg config.Config) error {
	return config.Save(&cfg)
}

// VerifyCloudPassword checks password with the specified cloud target URL
func (a *App) VerifyCloudPassword(targetURL string, password string) (bool, error) {
	if targetURL == "" {
		targetURL = a.getAPIURL()
	}
	url := fmt.Sprintf("%s/verify-cloud-access", targetURL)
	payload := map[string]string{"password": password}
	data, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(a.ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("Unable to connect to Cloud API server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	var res struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err == nil && res.Message != "" {
		return false, fmt.Errorf("%s", res.Message)
	}

	return false, fmt.Errorf("Access denied (status: %d)", resp.StatusCode)
}
