package apikey

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Service interface for API key operations
type Service interface {
	// GenerateAPIKey creates a new API key for a user
	GenerateAPIKey(userID uint, name string, expiry *time.Time, permissions []string) (string, *APIKey, error)
	
	// ValidateAPIKey checks if an API key is valid
	ValidateAPIKey(apiKey string) (*APIKey, error)
	
	// GetAPIKey gets an API key by ID
	GetAPIKey(id uint) (*APIKey, error)
	
	// ListAPIKeys lists all API keys for a user with pagination
	ListAPIKeys(userID uint, page, pageSize int) ([]*APIKey, int64, error)
	
	// RevokeAPIKey revokes (deletes) an API key
	RevokeAPIKey(id uint, userID uint) error
	
	// UpdateAPIKey updates an API key's name, permissions or expiry
	UpdateAPIKey(id uint, userID uint, name string, expiry *time.Time, permissions []string) (*APIKey, error)
}

// service is the implementation of Service interface
type service struct {
	repository Repository
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(repository Repository) Service {
	return &service{repository: repository}
}

// GenerateAPIKey creates a new API key for a user
func (s *service) GenerateAPIKey(userID uint, name string, expiry *time.Time, permissions []string) (string, *APIKey, error) {
	// Generate a random API key (32 bytes, 64 hex chars)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}
	
	keyString := hex.EncodeToString(b)
	
	// Get prefix for easy identification
	prefix := keyString[:8]
	
	// Hash the key for storage
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(keyString), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	
	// Convert permissions array to string
	permissionsStr := strings.Join(permissions, ",")
	
	apiKey := &APIKey{
		Name:        name,
		Key:         string(hashedKey),
		Prefix:      prefix,
		UserID:      userID,
		ExpiresAt:   expiry,
		Permissions: permissionsStr,
	}
	
	// Save to database
	if err := s.repository.Create(apiKey); err != nil {
		return "", nil, err
	}
	
	// Return the full key (will only be shown once to the user)
	return keyString, apiKey, nil
}

// ValidateAPIKey checks if an API key is valid and returns the API key entity
func (s *service) ValidateAPIKey(apiKeyString string) (*APIKey, error) {
	if len(apiKeyString) < 8 {
		return nil, errors.New("invalid API key format")
	}
	
	// Extract prefix (first 8 chars)
	prefix := apiKeyString[:8]
	
	// Find the API key by prefix
	apiKey, err := s.repository.FindByPrefix(prefix)
	if err != nil {
		return nil, errors.New("invalid API key")
	}
	
	// Check if key is expired
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("API key expired")
	}
	
	// Verify the key
	if err := bcrypt.CompareHashAndPassword([]byte(apiKey.Key), []byte(apiKeyString)); err != nil {
		return nil, errors.New("invalid API key")
	}
	
	// Update last used timestamp
	if err := s.repository.UpdateLastUsed(apiKey.ID); err != nil {
		// Non-critical error, just log it
		// logger.Warn("Failed to update API key last used timestamp", err)
	}
	
	return apiKey, nil
}

// GetAPIKey gets an API key by ID
func (s *service) GetAPIKey(id uint) (*APIKey, error) {
	return s.repository.FindByID(id)
}

// ListAPIKeys lists all API keys for a user with pagination
func (s *service) ListAPIKeys(userID uint, page, pageSize int) ([]*APIKey, int64, error) {
	return s.repository.FindByUserID(userID, page, pageSize)
}

// RevokeAPIKey revokes (deletes) an API key
func (s *service) RevokeAPIKey(id uint, userID uint) error {
	apiKey, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	
	// Security check: ensure the key belongs to the user
	if apiKey.UserID != userID {
		return errors.New("unauthorized to revoke this API key")
	}
	
	return s.repository.Delete(id)
}

// UpdateAPIKey updates an API key's name, permissions or expiry
func (s *service) UpdateAPIKey(id uint, userID uint, name string, expiry *time.Time, permissions []string) (*APIKey, error) {
	apiKey, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Security check: ensure the key belongs to the user
	if apiKey.UserID != userID {
		return nil, errors.New("unauthorized to update this API key")
	}
	
	// Update fields
	apiKey.Name = name
	apiKey.ExpiresAt = expiry
	apiKey.Permissions = strings.Join(permissions, ",")
	
	if err := s.repository.Update(apiKey); err != nil {
		return nil, err
	}
	
	return apiKey, nil
}
