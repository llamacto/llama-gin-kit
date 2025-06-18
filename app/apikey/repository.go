package apikey

import (
	"time"

	"gorm.io/gorm"
)

// Repository interface for API key operations
type Repository interface {
	Create(apiKey *APIKey) error
	FindByID(id uint) (*APIKey, error)
	FindByKey(key string) (*APIKey, error)
	FindByPrefix(prefix string) (*APIKey, error)
	FindByUserID(userID uint, page, pageSize int) ([]*APIKey, int64, error)
	Update(apiKey *APIKey) error
	Delete(id uint) error
	UpdateLastUsed(id uint) error
}

// repository is the implementation of Repository interface
type repository struct {
	db *gorm.DB
}

// NewAPIKeyRepository creates a new API key repository
func NewAPIKeyRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new API key
func (r *repository) Create(apiKey *APIKey) error {
	return r.db.Create(apiKey).Error
}

// FindByID finds an API key by its ID
func (r *repository) FindByID(id uint) (*APIKey, error) {
	var apiKey APIKey
	if err := r.db.First(&apiKey, id).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// FindByKey finds an API key by its key
func (r *repository) FindByKey(key string) (*APIKey, error) {
	var apiKey APIKey
	if err := r.db.Where("key = ?", key).First(&apiKey).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// FindByPrefix finds an API key by its prefix
func (r *repository) FindByPrefix(prefix string) (*APIKey, error) {
	var apiKey APIKey
	if err := r.db.Where("prefix = ?", prefix).First(&apiKey).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// FindByUserID finds all API keys for a user with pagination
func (r *repository) FindByUserID(userID uint, page, pageSize int) ([]*APIKey, int64, error) {
	var apiKeys []*APIKey
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := r.db.Model(&APIKey{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&apiKeys).Error; err != nil {
		return nil, 0, err
	}

	return apiKeys, total, nil
}

// Update updates an API key
func (r *repository) Update(apiKey *APIKey) error {
	return r.db.Save(apiKey).Error
}

// Delete soft deletes an API key
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&APIKey{}, id).Error
}

// UpdateLastUsed updates the last used timestamp for an API key
func (r *repository) UpdateLastUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&APIKey{}).Where("id = ?", id).Update("last_used_at", now).Error
}
