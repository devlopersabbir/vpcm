package config

import (
	"os"
	"path/filepath"
	"strings"
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
}

type APIConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
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

// Load loads config from default file, overrides from environment and/or flags
func Load() (*Config, error) {
	// Defaults
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.path", filepath.Join(os.Getenv("HOME"), ".local/share/vpsm/vpsm.db"))
	viper.SetDefault("api.enabled", false)
	viper.SetDefault("api.host", "127.0.0.1")
	viper.SetDefault("api.port", 8080)
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
