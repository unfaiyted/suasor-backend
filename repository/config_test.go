package repository_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"suasor/repository"
	"suasor/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestConfiguration creates a sample configuration for testing
func createTestConfiguration() *types.Configuration {
	// Create a configuration with test values
	// Update this to match your actual Configuration struct
	return &types.Configuration{}
}

func TestNewConfigRepository(t *testing.T) {
	t.Run("DefaultPath", func(t *testing.T) {
		os.Unsetenv("SUASOR_CONFIG_DIR")

		repo := repository.NewConfigRepository()
		assert.NotNil(t, repo)
	})

	t.Run("CustomPath", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "config-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		customPath := filepath.Join(tempDir, "custom-config.json")
		os.Setenv("SUASOR_CONFIG_DIR", customPath)
		defer os.Unsetenv("SUASOR_CONFIG_DIR")

		repo := repository.NewConfigRepository()
		assert.NotNil(t, repo)
	})
}

func TestEnsureConfigDir(t *testing.T) {
	// Save current directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	// Verify config dir doesn't exist yet
	_, err = os.Stat("./config")
	assert.True(t, os.IsNotExist(err))

	// Run the function
	repo := repository.NewConfigRepository()
	err = repo.EnsureConfigDir()
	assert.NoError(t, err)

	// Verify directory was created
	info, err := os.Stat("./config")
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestReadWriteConfigFile(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	os.Setenv("SUASOR_CONFIG_DIR", configPath)
	defer os.Unsetenv("SUASOR_CONFIG_DIR")

	// Create repository
	repo := repository.NewConfigRepository()

	// Create test config
	testConfig := createTestConfiguration()

	// Write config
	err = repo.WriteConfigFile(testConfig)
	require.NoError(t, err)

	// Read config back
	readConfig, err := repo.ReadConfigFile()
	require.NoError(t, err)

	// Compare configs - add specific field comparisons based on your Configuration struct
	assert.NotNil(t, readConfig)
}

func TestReadConfigFileNonExistent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	nonExistentPath := filepath.Join(tempDir, "non-existent.json")
	os.Setenv("SUASOR_CONFIG_DIR", nonExistentPath)
	defer os.Unsetenv("SUASOR_CONFIG_DIR")

	repo := repository.NewConfigRepository()
	config, err := repo.ReadConfigFile()

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "error reading config file")
}

func TestWatchConfigFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping watch test in short mode")
	}

	tempDir, err := os.MkdirTemp("", "config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	os.Setenv("SUASOR_CONFIG_DIR", configPath)
	defer os.Unsetenv("SUASOR_CONFIG_DIR")

	// Create initial config
	repo := repository.NewConfigRepository()
	testConfig := createTestConfiguration()
	err = repo.WriteConfigFile(testConfig)
	require.NoError(t, err)

	// Set up change detection
	var wg sync.WaitGroup
	wg.Add(1)
	changeDetected := false

	// Set up watch
	err = repo.WatchConfigFile(func() {
		changeDetected = true
		wg.Done()
	})
	assert.NoError(t, err)

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify the config
	err = repo.WriteConfigFile(testConfig)
	assert.NoError(t, err)

	// Wait for change with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		assert.True(t, changeDetected)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for config change detection")
	}
}
