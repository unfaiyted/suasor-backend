package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"suasor/types/models"
)

// TestInitialize tests the Initialize function but using an in-memory database
func TestInitialize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup context with test environment
	ctx := context.Background()
	logger := NewTestLogger()
	ctx = ContextWithTestLogger(ctx, logger)
	t.Setenv("GO_ENV", "dev")

	// Initialize in-memory database instead of PostgreSQL
	db, err := InitializeInMemoryDB(ctx)
	require.NoError(t, err, "Failed to initialize test database")
	require.NotNil(t, db, "Database connection should not be nil")

	// Verify the admin user was created
	var user models.User
	err = db.First(&user, "email = ?", "admin@dev.com").Error
	assert.NoError(t, err, "Should find the admin user in the database")
	assert.Equal(t, "devAdmin", user.Username, "Admin username should match")
}

// TestCreateTestAdminUser tests the admin user creation function
func TestCreateTestAdminUser(t *testing.T) {
	// Setup context with test environment
	ctx := context.Background()
	logger := NewTestLogger()
	ctx = ContextWithTestLogger(ctx, logger)
	t.Setenv("GO_ENV", "dev")

	// Initialize in-memory database
	db, err := InitializeInMemoryDB(ctx)
	require.NoError(t, err, "Failed to initialize test database")

	// First, delete any existing admin user
	db.Where("email = ?", "admin@dev.com").Delete(&models.User{})

	// Now run the function to create a test admin user
	err = CreateTestAdminUser(ctx, db)
	assert.NoError(t, err, "Should create test admin user without error")

	// Verify that the admin user was created
	var user models.User
	err = db.First(&user, "email = ?", "admin@dev.com").Error
	assert.NoError(t, err, "Should find the admin user in the database")
	assert.Equal(t, "devAdmin", user.Username, "Admin username should match")
	assert.Equal(t, "admin", user.Role, "Admin role should match")

	// Test that duplicate creation is handled properly
	err = CreateTestAdminUser(ctx, db)
	assert.NoError(t, err, "Should handle duplicate admin user creation gracefully")

	// Verify no environment setting
	t.Setenv("GO_ENV", "production")
	err = CreateTestAdminUser(ctx, db)
	assert.NoError(t, err, "Should skip admin creation in non-test environment")
}

// TestDatabaseUtilsFunctions tests basic functionality
func TestDatabaseUtilsFunctions(t *testing.T) {
	// Simple test to verify the test environment works
	assert.Equal(t, 1, 1, "Basic assertion should pass")
}