package repository_test

import (
	"context"
	"testing"

	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate the Client model
	err = db.AutoMigrate(
		&models.Client[types.ClientConfig]{},
		&models.Client[types.JellyfinConfig]{},
		&models.Client[types.SonarrConfig]{},
		&models.Client[types.RadarrConfig]{},
		&models.Client[types.LidarrConfig]{},
		&models.Client[types.EmbyConfig]{},
		&models.Client[types.JellyfinConfig]{},
		&models.Client[types.PlexConfig]{},
		&models.Client[types.SubsonicConfig]{},
	)

	require.NoError(t, err)

	return db
}

func createTestClient(t *testing.T, db *gorm.DB, userID uint64, clientType types.ClientType) *models.Client[types.JellyfinConfig] {
	client := models.Client[types.JellyfinConfig]{
		UserID:   userID,
		Name:     "Test Client",
		Category: types.ClientCategoryMedia,
		Config: models.ClientConfigWrapper[types.JellyfinConfig]{Data: types.JellyfinConfig{Enabled: true, BaseURL: "http://local", APIKey: "11",
			Username: "admin"},
		}}

	err := db.Create(&client).Error
	require.NoError(t, err)
	return &client
}

func TestClientRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewClientRepository[types.JellyfinConfig](db)
	ctx := context.Background()

	// Create a client
	clientConfig := types.JellyfinConfig{
		Enabled:  false,
		BaseURL:  "http://local",
		APIKey:   "11",
		Username: "admin",
	}

	client := models.Client[types.JellyfinConfig]{
		UserID:   1,
		Name:     "Test Client",
		Category: types.ClientCategoryMedia,
		Config:   models.ClientConfigWrapper[types.JellyfinConfig]{Data: clientConfig},
	}

	result, err := repo.Create(ctx, client)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotZero(t, result.ID)
	assert.Equal(t, client.Name, result.Name)
	assert.Equal(t, client.UserID, result.UserID)
	assert.Equal(t, client.Category, result.Category)
}

func TestClientRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewClientRepository[types.EmbyConfig](db)
	ctx := context.Background()

	// Create a client first
	originalClient := models.Client[types.EmbyConfig]{
		UserID:   1,
		Name:     "Original Name",
		Category: types.ClientCategoryMedia,
		Config:   models.ClientConfigWrapper[types.EmbyConfig]{Data: types.NewEmbyConfig()},
	}

	created, err := repo.Create(ctx, originalClient)
	require.NoError(t, err)

	// Update the client
	updatedClient := *created
	updatedClient.Name = "Updated Name"

	result, err := repo.Update(ctx, updatedClient)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", result.Name)

	// Verify in the database
	var retrieved models.Client[types.EmbyConfig]
	err = db.Where("id = ?", created.ID).First(&retrieved).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", retrieved.Name)
}

func TestClientRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewClientRepository[types.PlexConfig](db)
	ctx := context.Background()

	// Create a client
	originalClient := models.Client[types.PlexConfig]{
		UserID:   1,
		Name:     "Test Client",
		Category: types.ClientCategoryMedia,
		Config:   models.ClientConfigWrapper[types.PlexConfig]{Data: types.NewPlexConfig()},
	}

	created, err := repo.Create(ctx, originalClient)
	require.NoError(t, err)

	// Get by ID
	result, err := repo.GetByID(ctx, created.ID, created.UserID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, created.Name, result.Name)

	// Test with invalid ID
	_, err = repo.GetByID(ctx, 9999, created.UserID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestClientRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewClientRepository[types.SubsonicConfig](db)
	ctx := context.Background()

	userID := uint64(1)
	otherUserID := uint64(2)

	// Create clients for our test user
	client1 := models.Client[types.SubsonicConfig]{
		UserID:   userID,
		Name:     "Client 1",
		Category: types.ClientCategoryMedia,
		Config:   models.ClientConfigWrapper[types.SubsonicConfig]{Data: types.NewSubsonicConfig()},
	}

	client2 := models.Client[types.SubsonicConfig]{
		UserID:   userID,
		Name:     "Client 2",
		Category: types.ClientCategoryMedia,
		Config:   models.ClientConfigWrapper[types.SubsonicConfig]{Data: types.NewSubsonicConfig()},
	}

	// Create a client for another user
	otherClient := models.Client[types.SubsonicConfig]{
		UserID:   otherUserID,
		Name:     "Other Client",
		Category: types.ClientCategoryMedia,
		Config:   models.ClientConfigWrapper[types.SubsonicConfig]{Data: types.NewSubsonicConfig()},
	}

	// Add clients to the database
	_, err := repo.Create(ctx, client1)
	require.NoError(t, err)

	_, err = repo.Create(ctx, client2)
	require.NoError(t, err)

	_, err = repo.Create(ctx, otherClient)
	require.NoError(t, err)

	// Get clients for our test user
	results, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify client names are in the results
	names := []string{results[0].Name, results[1].Name}
	assert.Contains(t, names, "Client 1")
	assert.Contains(t, names, "Client 2")

	// Get clients for the other user
	otherResults, err := repo.GetByUserID(ctx, otherUserID)
	require.NoError(t, err)
	assert.Len(t, otherResults, 1)
	assert.Equal(t, "Other Client", otherResults[0].Name)
}

func TestClientRepository_GetByCategory(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewClientRepository[types.RadarrConfig](db)
	ctx := context.Background()

	userID := uint64(1)

	// Create different types of clients
	radarrClient := models.Client[types.RadarrConfig]{
		UserID:   userID,
		Name:     "Radarr Client",
		Category: types.ClientCategoryAutomation,
		Config:   models.ClientConfigWrapper[types.RadarrConfig]{Data: types.NewRadarrConfig()},
	}

	radarrClient2 := models.Client[types.RadarrConfig]{
		UserID:   userID,
		Name:     "Second Radarr Client",
		Category: types.ClientCategoryAutomation,
		Config:   models.ClientConfigWrapper[types.RadarrConfig]{Data: types.NewRadarrConfig()},
	}

	// Add clients to the database
	_, err := repo.Create(ctx, radarrClient)
	require.NoError(t, err)

	_, err = repo.Create(ctx, radarrClient2)
	require.NoError(t, err)

	// Get clients by type
	radarrResults, err := repo.GetByCategory(ctx, types.ClientCategoryAutomation, userID)
	require.NoError(t, err)
	assert.Len(t, radarrResults, 2)
	assert.Equal(t, "Radarr Client", radarrResults[0].Name)

	assert.Len(t, radarrResults, 2)
	assert.Equal(t, "Second Radarr Client", radarrResults[1].Name)
}

func TestClientRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewClientRepository[types.EmbyConfig](db)

	ctx := context.Background()

	// Create a client
	client := models.Client[types.EmbyConfig]{
		UserID:   1,
		Name:     "Test Client",
		Category: types.ClientCategoryMedia,
		Config:   models.ClientConfigWrapper[types.EmbyConfig]{Data: types.NewEmbyConfig()},
	}

	created, err := repo.Create(ctx, client)
	require.NoError(t, err)

	// Delete the client
	err = repo.Delete(ctx, created.ID, created.UserID)
	require.NoError(t, err)

	// Verify it's gone
	_, err = repo.GetByID(ctx, created.ID, created.UserID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	// Try to delete a non-existent client
	err = repo.Delete(ctx, 9999, created.UserID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
