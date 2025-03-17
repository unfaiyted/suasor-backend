// repository/config.go
package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"suasor/models"

	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// ConfigRepository handles configuration storage operations
type ConfigRepository interface {
	ReadConfigFile() (*models.Configuration, error)
	WriteConfigFile(cfg *models.Configuration) error
	WatchConfigFile(onChange func()) error
	EnsureConfigDir() error
}

type configRepository struct {
	configPath string
}

// NewConfigRepository creates a new configuration repository
func NewConfigRepository() ConfigRepository {
	if configPath := os.Getenv("SUASOR_CONFIG_DIR"); configPath != "" {
		return &configRepository{
			configPath: configPath,
		}
	} else {
		return &configRepository{
			configPath: "./config/app.config.json",
		}
	}
}

// EnsureConfigDir ensures the configuration directory exists
func (r *configRepository) EnsureConfigDir() error {
	dir := "./config"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}
	return nil
}

// ReadConfigFile reads the configuration file
func (r *configRepository) ReadConfigFile() (*models.Configuration, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(r.configPath), kjson.Parser()); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &models.Configuration{}
	if err := k.Unmarshal("", config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

// WriteConfigFile writes the configuration to file
func (r *configRepository) WriteConfigFile(cfg *models.Configuration) error {
	// Convert config struct to map
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &configMap); err != nil {
		return fmt.Errorf("error unmarshaling to map: %w", err)
	}

	// Save to file
	data, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile(r.configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// WatchConfigFile sets up a watcher for the config file
func (r *configRepository) WatchConfigFile(onChange func()) error {
	fp := file.Provider(r.configPath)
	fp.Watch(func(event interface{}, err error) {
		if err != nil {
			fmt.Printf("watch error: %v\n", err)
			return
		}
		onChange()
	})
	return nil
}
