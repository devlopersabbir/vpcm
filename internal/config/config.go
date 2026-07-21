package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database  DatabaseConfig  `mapstructure:"database"`
	API       APIConfig       `mapstructure:"api"`
	SSH       SSHConfig       `mapstructure:"ssh"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Collector CollectorConfig `mapstructure:"collector"`
	Plugins   PluginsConfig   `mapstructure:"plugins"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Path   string `mapstructure:"path"`
	URI    string `mapstructure:"uri"`
	Name   string `mapstructure:"name"`
}

type APIConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Mode      string `mapstructure:"mode"`       // "local" or "cloud"
	Token     string `mapstructure:"token"`      // cloud API auth token
	GlobalURL string `mapstructure:"global_url"` // SaaS API base URL
}

type SSHConfig struct {
	Timeout time.Duration `mapstructure:"timeout"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json or pretty
}

type CollectorConfig struct {
	Workers int `mapstructure:"workers"`
}

type PluginsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

var CConfig Config
var configMu sync.Mutex

// Load loads config from default file, overrides from environment and/or flags
func Load() (*Config, error) {
	configMu.Lock()
	defer configMu.Unlock()

	// Defaults
	viper.SetDefault("database.driver", "mongodb")
	viper.SetDefault("database.path", filepath.Join(os.Getenv("HOME"), ".local", "share", "vpsm", "vpsm.db"))
	viper.SetDefault("database.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.name", "vpsm")
	viper.SetDefault("api.enabled", false)
	viper.SetDefault("api.host", "127.0.0.1")
	viper.SetDefault("api.port", 8080)
	viper.SetDefault("api.mode", "local")
	viper.SetDefault("api.token", "")
	viper.SetDefault("api.global_url", "https://api.vpsm.dev")
	viper.SetDefault("ssh.timeout", 10*time.Second)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "pretty")
	viper.SetDefault("collector.workers", 5)
	viper.SetDefault("plugins.enabled", true)

	// File setup
	configHome := filepath.Join(os.Getenv("HOME"), ".config", "vpsm")
	viper.AddConfigPath(configHome)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read config file if present
	_ = viper.ReadInConfig()

	// Env setup
	viper.SetEnvPrefix("VPSM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Parse to struct
	if err := viper.Unmarshal(&CConfig); err != nil {
		return nil, err
	}

	// Expand home directory if sqlite path contains ~
	if strings.HasPrefix(CConfig.Database.Path, "~") {
		home := os.Getenv("HOME")
		CConfig.Database.Path = filepath.Join(home, CConfig.Database.Path[1:])
	}

	return &CConfig, nil
}

// Save writes the given config back to ~/.config/vpsm/config.yaml
func Save(cfg *Config) error {
	configHome := filepath.Join(os.Getenv("HOME"), ".config", "vpsm")
	if err := os.MkdirAll(configHome, 0755); err != nil {
		return err
	}
	configPath := filepath.Join(configHome, "config.yaml")

	viper.Set("database.driver", cfg.Database.Driver)
	viper.Set("database.path", cfg.Database.Path)
	viper.Set("database.uri", cfg.Database.URI)
	viper.Set("database.name", cfg.Database.Name)
	viper.Set("api.enabled", cfg.API.Enabled)
	viper.Set("api.host", cfg.API.Host)
	viper.Set("api.port", cfg.API.Port)
	viper.Set("api.mode", cfg.API.Mode)
	viper.Set("api.token", cfg.API.Token)
	viper.Set("api.global_url", cfg.API.GlobalURL)
	viper.Set("ssh.timeout", cfg.SSH.Timeout)
	viper.Set("logging.level", cfg.Logging.Level)
	viper.Set("logging.format", cfg.Logging.Format)
	viper.Set("collector.workers", cfg.Collector.Workers)
	viper.Set("plugins.enabled", cfg.Plugins.Enabled)

	return viper.WriteConfigAs(configPath)
}
