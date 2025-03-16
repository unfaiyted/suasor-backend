package repository

import (
	"context"
	"suasor/models"

	"gorm.io/gorm"
)

// ShortenRepository defines the data access interface for URL shortening
type ShortenRepository interface {
	GetAll(ctx context.Context) ([]models.Shorten, error)
	FindById(ctx context.Context, id uint64) (*models.Shorten, error)
	FindByCode(ctx context.Context, code string) (*models.Shorten, error)
	Create(ctx context.Context, shorten *models.Shorten) (*models.Shorten, error)
	Update(ctx context.Context, shorten *models.Shorten) (*models.Shorten, error)
	Delete(ctx context.Context, code string) (string, error)
	IncrementClickCount(ctx context.Context, code string) (*models.Shorten, error)
	FindByOriginalURL(ctx context.Context, url string) (*models.Shorten, error)
}

type shortenRepository struct {
	db *gorm.DB
}

// NewConfigRepository creates a new configuration repository
func NewShortenRepository(db *gorm.DB) ShortenRepository {
	return &shortenRepository{
		db: db,
	}
}

func (r *shortenRepository) GetAll(ctx context.Context) ([]models.Shorten, error) {
	var shortens []models.Shorten

	result := r.db.Find(&shortens)
	return shortens, result.Error
}

func (r *shortenRepository) FindById(ctx context.Context, id uint64) (*models.Shorten, error) {
	var shorten models.Shorten
	result := r.db.First(&shorten, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &shorten, result.Error
}

func (r *shortenRepository) FindByCode(ctx context.Context, code string) (*models.Shorten, error) {
	var shorten models.Shorten
	result := r.db.Where("short_code", code).First(&shorten)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &shorten, result.Error
}

func (r *shortenRepository) Create(ctx context.Context, shorten *models.Shorten) (*models.Shorten, error) {
	result := r.db.Create(&shorten)
	return shorten, result.Error
}

func (r *shortenRepository) Update(ctx context.Context, shorten *models.Shorten) (*models.Shorten, error) {
	result := r.db.Save(&shorten)
	return shorten, result.Error
}

func (r *shortenRepository) IncrementClickCount(ctx context.Context, code string) (*models.Shorten, error) {
	var shorten models.Shorten
	r.db.Where("short_code", code).First(&shorten)

	shorten.ClickCount = shorten.ClickCount + 1

	result := r.db.Save(&shorten)
	return &shorten, result.Error
}

func (r *shortenRepository) Delete(ctx context.Context, code string) (string, error) {
	var shorten models.Shorten
	result := r.db.Delete(&shorten, code)
	return code, result.Error
}

func (r *shortenRepository) FindByOriginalURL(ctx context.Context, url string) (*models.Shorten, error) {
	var shorten models.Shorten
	result := r.db.Where("original_url = ?", url).First(&shorten)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &shorten, nil
}
