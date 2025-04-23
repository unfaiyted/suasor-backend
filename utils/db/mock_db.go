package database

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"suasor/clients/media/types"
	dbTypes "suasor/types"
	"suasor/types/models"
)

// MockDB is a mock implementation of the database for unit tests
type MockDB struct {
	mock.Mock
}

// MockGormDB creates a mock DB that can be used as a replacement for *gorm.DB in tests
type MockGormDB struct {
	mock *MockDB
}

// Initialize mocks the database initialization
func (m *MockDB) Initialize(ctx context.Context, config dbTypes.DatabaseConfig) (*gorm.DB, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// CreateTestAdminUser mocks the admin user creation function
func (m *MockDB) CreateTestAdminUser(ctx context.Context, db *gorm.DB) error {
	args := m.Called(ctx, db)
	return args.Error(0)
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{}
}

// Where is a mock implementation for the Where method in GORM
func (m *MockGormDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.mock.Called(query, args)
	return &gorm.DB{} // Return a mock gorm.DB
}

// First is a mock implementation for the First method in GORM
func (m *MockGormDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	m.mock.Called(dest, conds)
	return &gorm.DB{}
}

// Create is a mock implementation for the Create method in GORM
func (m *MockGormDB) Create(value interface{}) *gorm.DB {
	m.mock.Called(value)
	return &gorm.DB{}
}

// Save is a mock implementation for the Save method in GORM
func (m *MockGormDB) Save(value interface{}) *gorm.DB {
	m.mock.Called(value)
	return &gorm.DB{}
}

// Delete is a mock implementation for the Delete method in GORM
func (m *MockGormDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	m.mock.Called(value, conds)
	return &gorm.DB{}
}

// Find is a mock implementation for the Find method in GORM
func (m *MockGormDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	m.mock.Called(dest, conds)
	return &gorm.DB{}
}

// Count is a mock implementation for the Count method in GORM
func (m *MockGormDB) Count(count *int64) *gorm.DB {
	m.mock.Called(count)
	return &gorm.DB{}
}

// Preload is a mock implementation for the Preload method in GORM
func (m *MockGormDB) Preload(query string, args ...interface{}) *gorm.DB {
	m.mock.Called(query, args)
	return &gorm.DB{}
}

// Joins is a mock implementation for the Joins method in GORM
func (m *MockGormDB) Joins(query string, args ...interface{}) *gorm.DB {
	m.mock.Called(query, args)
	return &gorm.DB{}
}

// Begin is a mock implementation for the Begin method in GORM
func (m *MockGormDB) Begin() *gorm.DB {
	m.mock.Called()
	return &gorm.DB{}
}

// Commit is a mock implementation for the Commit method in GORM
func (m *MockGormDB) Commit() *gorm.DB {
	m.mock.Called()
	return &gorm.DB{}
}

// Rollback is a mock implementation for the Rollback method in GORM
func (m *MockGormDB) Rollback() *gorm.DB {
	m.mock.Called()
	return &gorm.DB{}
}

// Exec is a mock implementation for the Exec method in GORM
func (m *MockGormDB) Exec(sql string, values ...interface{}) *gorm.DB {
	m.mock.Called(sql, values)
	return &gorm.DB{}
}

// Raw is a mock implementation for the Raw method in GORM
func (m *MockGormDB) Raw(sql string, values ...interface{}) *gorm.DB {
	m.mock.Called(sql, values)
	return &gorm.DB{}
}

// Order is a mock implementation for the Order method in GORM
func (m *MockGormDB) Order(value interface{}) *gorm.DB {
	m.mock.Called(value)
	return &gorm.DB{}
}

// Limit is a mock implementation for the Limit method in GORM
func (m *MockGormDB) Limit(limit int) *gorm.DB {
	m.mock.Called(limit)
	return &gorm.DB{}
}

// Offset is a mock implementation for the Offset method in GORM
func (m *MockGormDB) Offset(offset int) *gorm.DB {
	m.mock.Called(offset)
	return &gorm.DB{}
}

// Table is a mock implementation for the Table method in GORM
func (m *MockGormDB) Table(name string) *gorm.DB {
	m.mock.Called(name)
	return &gorm.DB{}
}

// Update is a mock implementation for the Update method in GORM
func (m *MockGormDB) Update(column string, value interface{}) *gorm.DB {
	m.mock.Called(column, value)
	return &gorm.DB{}
}

// Updates is a mock implementation for the Updates method in GORM
func (m *MockGormDB) Updates(values interface{}) *gorm.DB {
	m.mock.Called(values)
	return &gorm.DB{}
}

// Model is a mock implementation for the Model method in GORM
func (m *MockGormDB) Model(value interface{}) *gorm.DB {
	m.mock.Called(value)
	return &gorm.DB{}
}

// WithContext is a mock implementation for the WithContext method in GORM
func (m *MockGormDB) WithContext(ctx context.Context) *gorm.DB {
	m.mock.Called(ctx)
	return &gorm.DB{}
}

// Transaction is a mock implementation for the Transaction method in GORM
func (m *MockGormDB) Transaction(fc func(tx *gorm.DB) error) error {
	args := m.mock.Called(fc)
	return args.Error(0)
}

// MockUser creates a mock user for testing
func MockUser(id uint64, username, email, role string) models.User {
	user := models.User{
		Username: username,
		Email:    email,
		Role:     role,
	}
	user.ID = id
	return user
}

// MockMediaItem creates a mock media item for testing
func MockMediaItem[T types.MediaData](id uint64, title string, mediaType types.MediaType, data T) models.MediaItem[T] {
	now := time.Now()
	item := models.MediaItem[T]{
		ID:          id,
		Type:        mediaType,
		Title:       title,
		Data:        data,
		ReleaseYear: 2023,
		ReleaseDate: now.AddDate(-1, 0, 0),
	}
	item.CreatedAt = now.Add(-7 * 24 * time.Hour)
	item.UpdatedAt = now
	return item
}

// MockMovie creates a mock movie for testing
func MockMovie(id uint64, title string) *models.MediaItem[*types.Movie] {
	movieData := &types.Movie{
		Details: types.MediaDetails{
			Title:       title,
			Description: "A mock movie for testing",
			ReleaseYear: 2023,
		},
	}
	movie := MockMediaItem(id, title, types.MediaTypeMovie, movieData)
	return &movie
}

// MockSeries creates a mock series for testing
func MockSeries(id uint64, title string) *models.MediaItem[*types.Series] {
	seriesData := &types.Series{
		Details: types.MediaDetails{
			Title:       title,
			Description: "A mock series for testing",
			ReleaseYear: 2023,
		},
	}
	series := MockMediaItem(id, title, types.MediaTypeSeries, seriesData)
	return &series
}

// MockUserMediaItemData creates mock user media item data for testing
func MockUserMediaItemData[T types.MediaData](id uint64, userID uint64, mediaItemID uint64, mediaType types.MediaType) models.UserMediaItemData[T] {
	now := time.Now()
	data := models.UserMediaItemData[T]{
		ID:              id,
		UserID:          userID,
		MediaItemID:     mediaItemID,
		Type:            mediaType,
		PlayedAt:        now.Add(-24 * time.Hour),
		LastPlayedAt:    now,
		PlayCount:       1,
		IsFavorite:      true,
		UserRating:      8.5,
		Completed:       false,
		PositionSeconds: 3600,
		DurationSeconds: 7200,
		CreatedAt:       now.Add(-7 * 24 * time.Hour),
		UpdatedAt:       now,
	}
	return data
}

// SetupMockDB configures a mock DB with common expectations
func SetupMockDB(t mock.TestingT) *MockDB {
	mockDB := NewMockDB()
	return mockDB
}