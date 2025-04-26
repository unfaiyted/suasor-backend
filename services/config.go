package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"suasor/repository"
	"suasor/types"
	"suasor/types/constants"
	"suasor/utils/logger"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/dotenv"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// tryLock attempts to acquire a lock with a timeout
func tryLock(lock *sync.RWMutex) bool {
	// Create a channel to signal when the lock is acquired
	done := make(chan bool, 1)
	
	// Try to acquire the lock in a goroutine
	go func() {
		lock.Lock()
		done <- true
	}()
	
	// Wait for the lock with timeout
	select {
	case <-done:
		// Lock acquired
		return true
	case <-time.After(500 * time.Millisecond):
		// Timeout - couldn't acquire lock
		return false
	}
}

// ConfigService provides methods to interact with configuration
type ConfigService interface {
	InitConfig(ctx context.Context) error
	GetConfig() *types.Configuration
	SaveConfig(ctx context.Context, cfg types.Configuration) error
	GetFileConfig(ctx context.Context) *types.Configuration
	SaveFileConfig(ctx context.Context, cfg types.Configuration) error
	ResetFileConfig(ctx context.Context) error
	GetRepo() repository.ConfigRepository
}

type configService struct {
	configRepo repository.ConfigRepository
	config     *types.Configuration
	configLock sync.RWMutex
	k          *koanf.Koanf
	configPath string
}

// NewConfigService creates a new configuration service
func NewConfigService(configRepo repository.ConfigRepository) ConfigService {
	return &configService{
		configRepo: configRepo,
		configPath: "config/app.config.json",
	}
}

// InitConfig initializes the configuration
func (s *configService) InitConfig(ctx context.Context) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Initializing Config")
	s.k = koanf.New(".")

	if err := godotenv.Load(); err != nil {
		log.Info().Msg("Unable to load env")
	}

	// 1. Load defaults
	log.Debug().Msg("Loading default configuration")
	if err := s.k.Load(confmap.Provider(constants.DefaultConfig, "."), nil); err != nil {
		log.Error().Err(err).Msg("Error loading defaults")
		return fmt.Errorf("error loading defaults: %w", err)
	}
	log.Debug().Interface("defaults", constants.DefaultConfig).Msg("Default configuration loaded")

	// Ensure config directory exists
	log.Debug().Msg("Ensuring config directory exists")
	if err := s.configRepo.EnsureConfigDir(); err != nil {
		log.Error().Err(err).Msg("Failed to ensure config directory")
		return err
	}

	// 2. Load app.config.json
	log.Debug().Str("path", s.configPath).Msg("Loading configuration from file")
	f := file.Provider(s.configPath)

	if err := s.k.Load(f, kjson.Parser()); err != nil {
		// Create default config if file doesn't exist
		if os.IsNotExist(err) {
			log.Info().Msg("Config file doesn't exist, creating default")
			defaultConfig := &types.Configuration{}
			if err := s.k.Unmarshal("", defaultConfig); err != nil {
				log.Error().Err(err).Msg("Error unmarshaling default config")
				return fmt.Errorf("error unmarshaling default config: %w", err)
			}
			if err := s.configRepo.WriteConfigFile(defaultConfig); err != nil {
				log.Error().Err(err).Msg("Error saving default config")
				return fmt.Errorf("error saving default config: %w", err)
			}
			log.Info().Msg("Default config file created successfully")
		} else {
			log.Error().Err(err).Msg("Error loading config file")
			return fmt.Errorf("error loading config file: %w", err)
		}
	} else {
		log.Info().Msg("Config file loaded successfully")
	}

	// Set up file watcher
	log.Debug().Msg("Setting up config file watcher")
	s.configRepo.WatchConfigFile(func() {
		// Use a separate goroutine to prevent blocking
		go func() {
			logFromWatcher := log.With().Str("source", "file_watcher").Logger()
			logFromWatcher.Info().Msg("Config file change detected")

			// Try to acquire the lock, but don't block indefinitely
			// if we can't get the lock, skip this reload
			if !tryLock(&s.configLock) {
				logFromWatcher.Warn().Msg("Config is being modified elsewhere, skipping reload")
				return
			}
			defer s.configLock.Unlock()

			// Create a new koanf instance to avoid concurrent map issues
			newK := koanf.New(".")

			// Reload in the correct order
			logFromWatcher.Debug().Msg("Reloading default configuration")
			if err := newK.Load(confmap.Provider(constants.DefaultConfig, "."), nil); err != nil {
				logFromWatcher.Error().Err(err).Msg("Failed to reload default configuration")
				return
			}

			logFromWatcher.Debug().Msg("Reloading configuration from file")
			if err := newK.Load(file.Provider(s.configPath), kjson.Parser()); err != nil {
				logFromWatcher.Error().Err(err).Msg("Failed to reload configuration from file")
				return
			}

			logFromWatcher.Debug().Msg("Reloading configuration from .env file")
			newK.Load(file.Provider(".env"), dotenv.Parser())

			logFromWatcher.Debug().Msg("Reloading configuration from environment variables")
			newK.Load(env.Provider("suasor_", ".", s.envKeyReplacer), nil)

			logFromWatcher.Info().Msg("All config providers reloaded")

			// Create a new config struct to avoid modifying the existing one during unmarshaling
			newConfig := &types.Configuration{}
			if err := newK.Unmarshal("", newConfig); err != nil {
				logFromWatcher.Error().Err(err).Msg("Error unmarshaling configuration")
				return
			}

			// Now update the service's config and koanf instance
			s.k = newK
			s.config = newConfig

			logFromWatcher.Info().Msg("Configuration reloaded successfully due to file change")
		}()
	})

	// 3. Load environment variables
	log.Debug().Msg("Loading configuration from environment variables")
	envVars := make(map[string]string)
	log.Debug().Strs("envs", os.Environ()).Msg("envs")
	for _, e := range os.Environ() {

		if strings.HasPrefix(e, "suasor_") {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				log.Debug().Str("key", key).Str("value", value).Msg("Key values")
				transformedKey := s.envKeyReplacer(key)
				envVars[key] = value
				log.Debug().Str("original", key).Str("transformed", transformedKey).Str("value", value).Msg("Processing env var")
			}
		}
	}
	log.Debug().Interface("env_vars", envVars).Msg("Environment variables detected")

	// Use a custom env provider that can handle arrays
	koanfEnv := env.ProviderWithValue("suasor_", ".", func(key, value string) (string, interface{}) {
		// First apply the standard key transformation
		k := s.envKeyReplacer(key)

		// Check if this key should be an array
		if k == "auth.allowedOrigins" {
			log.Debug().Str("key", k).Str("value", value).Msg("Parsing array from env var")
			// Split by comma and trim whitespace
			parts := strings.Split(value, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			return k, parts
		}

		// For regular values, return as-is
		return k, value
	})

	if err := s.k.Load(koanfEnv, nil); err != nil {
		log.Error().Err(err).Msg("Error loading environment variables")
		return fmt.Errorf("error loading environment variables: %w", err)
	}

	// Load the final config
	s.configLock.Lock()
	defer s.configLock.Unlock()

	// Log the final merged configuration for testing
	var rawMergedConfig map[string]interface{}
	if err := s.k.Unmarshal("", &rawMergedConfig); err == nil {
		log.Debug().Interface("final_merged_config", rawMergedConfig).Msg("Final merged configuration")
	}

	s.config = &types.Configuration{}
	if err := s.k.UnmarshalWithConf("", s.config, koanf.UnmarshalConf{
		Tag: "json",
	}); err != nil {
		log.Error().Err(err).Msg("Error unmarshaling config")
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	log.Info().Msg("Configuration initialized successfully")
	return nil
}

// Helper method for environment variable key conversioenvsn
func (s *configService) envKeyReplacer(key string) string {

	// original := key
	transformed := strings.ReplaceAll(
		strings.ToLower(
			strings.TrimPrefix(key, "suasor_")),
		"_",
		".",
	)

	// TODO: auto detect the matching keys
	if transformed == "auth.allowedorigins" {
		transformed = "auth.allowedOrigins"
	}
	if transformed == "app.appurl" {
		transformed = "app.appURL"
	}

	// This would normally use logger, but since this is called during config loading
	// before logger might be fully available, we don't log here
	return transformed
}

// GetConfig returns the current configuration
func (s *configService) GetConfig() *types.Configuration {
	s.configLock.RLock()
	defer s.configLock.RUnlock()
	return s.config
}

// SaveConfig saves and updates the configuration
func (s *configService) SaveConfig(ctx context.Context, cfg types.Configuration) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Saving configuration")

	s.configLock.Lock()
	defer s.configLock.Unlock()

	// Convert config struct to map
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling config")
		return fmt.Errorf("error marshaling config: %w", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &configMap); err != nil {
		log.Error().Err(err).Msg("Error unmarshaling to map")
		return fmt.Errorf("error unmarshaling to map: %w", err)
	}

	log.Debug().Interface("config", configMap).Msg("New configuration to be applied")

	// Load the new config
	if err := s.k.Load(confmap.Provider(configMap, "."), nil); err != nil {
		log.Error().Err(err).Msg("Error loading new config")
		return fmt.Errorf("error loading new config: %w", err)
	}

	// Save to file
	log.Debug().Msg("Writing configuration to file")
	if err := s.configRepo.WriteConfigFile(&cfg); err != nil {
		log.Error().Err(err).Msg("Error saving config to file")
		return fmt.Errorf("error saving config: %w", err)
	}

	s.config = &cfg
	log.Info().Msg("Configuration saved successfully")
	return nil
}

// GetFileConfig returns only the file-based configuration
func (s *configService) GetFileConfig(ctx context.Context) *types.Configuration {
	log := logger.LoggerFromContext(ctx)
	log.Debug().Msg("Reading configuration from file")

	cfg, err := s.configRepo.ReadConfigFile()
	if err != nil {
		log.Error().Err(err).Msg("Error reading config file")
		return nil
	}

	log.Debug().Interface("file_config", cfg).Msg("File configuration retrieved")
	return cfg
}

// SaveFileConfig saves the configuration to file only
func (s *configService) SaveFileConfig(ctx context.Context, cfg types.Configuration) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Saving configuration to file only")

	err := s.configRepo.WriteConfigFile(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("Error writing config to file")
		return err
	}

	log.Info().Msg("Configuration file saved successfully")
	return nil
}

// ResetFileConfig resets config file to defaults
func (s *configService) ResetFileConfig(ctx context.Context) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Resetting configuration file to defaults")

	// Acquire lock to prevent concurrent access with file watcher
	s.configLock.Lock()
	defer s.configLock.Unlock()

	// Create default config
	k := koanf.New(".")
	log.Debug().Msg("Loading default configuration")
	if err := k.Load(confmap.Provider(constants.DefaultConfig, "."), nil); err != nil {
		log.Error().Err(err).Msg("Error loading defaults")
		return fmt.Errorf("error loading defaults: %w", err)
	}

	defaultConfig := &types.Configuration{}
	if err := k.Unmarshal("", defaultConfig); err != nil {
		log.Error().Err(err).Msg("Error unmarshaling default config")
		return fmt.Errorf("error unmarshaling default config: %w", err)
	}

	log.Debug().Interface("default_config", defaultConfig).Msg("Default configuration created")

	// Save defaults to file
	log.Debug().Msg("Writing default configuration to file")
	if err := s.configRepo.WriteConfigFile(defaultConfig); err != nil {
		log.Error().Err(err).Msg("Error writing default config to file")
		return fmt.Errorf("error writing default config: %w", err)
	}

	log.Info().Msg("Default configuration saved to file")

	// Update the in-memory configuration directly instead of calling InitConfig
	// to avoid race conditions with the file watcher
	s.config = defaultConfig
	
	// Recreate the koanf instance with default values
	s.k = koanf.New(".")
	if err := s.k.Load(confmap.Provider(constants.DefaultConfig, "."), nil); err != nil {
		log.Error().Err(err).Msg("Error reloading defaults")
		return fmt.Errorf("error reloading defaults: %w", err)
	}

	log.Info().Msg("Configuration reset to defaults successfully")
	return nil
}

func (s *configService) GetRepo() repository.ConfigRepository {
	return s.configRepo
}
