package apikey

import (
	"time"
)

// CreateRequest represents the request to create an API key
type CreateRequest struct {
	Name        string    `json:"name" binding:"required,max=100"`
	Permissions []string  `json:"permissions" binding:"omitempty"`
	ExpiresAt   time.Time `json:"expires_at" binding:"omitempty"`
	NeverExpire bool      `json:"never_expire" binding:"omitempty"`
}

// UpdateRequest represents the request to update an API key
type UpdateRequest struct {
	Name        string    `json:"name" binding:"omitempty,max=100"`
	Permissions []string  `json:"permissions" binding:"omitempty"`
	ExpiresAt   time.Time `json:"expires_at" binding:"omitempty"`
	NeverExpire bool      `json:"never_expire" binding:"omitempty"`
}

// Response represents the response format for API key operations
type Response struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Prefix      string     `json:"prefix"`
	Key         string     `json:"key,omitempty"` // Only included when creating a new key
	UserID      uint       `json:"user_id"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ListResponse represents the paginated response for listing API keys
type ListResponse struct {
	Total   int64      `json:"total"`
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
	Data    []Response `json:"data"`
}

// ToResponse converts an APIKey model to a Response DTO
func ToResponse(apiKey *APIKey, includeKey string) Response {
	var permissions []string
	if apiKey.Permissions != "" {
		permissions = splitPermissions(apiKey.Permissions)
	}

	return Response{
		ID:          apiKey.ID,
		Name:        apiKey.Name,
		Prefix:      apiKey.Prefix,
		Key:         includeKey,
		UserID:      apiKey.UserID,
		ExpiresAt:   apiKey.ExpiresAt,
		LastUsedAt:  apiKey.LastUsedAt,
		Permissions: permissions,
		CreatedAt:   apiKey.CreatedAt,
	}
}

// ToResponseList converts a slice of APIKey models to a slice of Response DTOs
func ToResponseList(apiKeys []*APIKey) []Response {
	responses := make([]Response, len(apiKeys))
	for i, apiKey := range apiKeys {
		responses[i] = ToResponse(apiKey, "")
	}
	return responses
}

// Helper function to split permission string
func splitPermissions(permissions string) []string {
	if permissions == "" {
		return []string{}
	}
	return splitCSV(permissions)
}

// Helper function to split CSV values
func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	split := make([]string, 0)
	for _, item := range splitString(s, ",") {
		if item != "" {
			split = append(split, item)
		}
	}
	return split
}

// Helper function to split string by separator
func splitString(s string, separator string) []string {
	if s == "" {
		return []string{}
	}
	return splitWithoutEmpty(s, separator)
}

// Helper function to split string and filter out empty values
func splitWithoutEmpty(s string, separator string) []string {
	result := make([]string, 0)
	for _, item := range split(s, separator) {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

// Helper function to split string by separator
func split(s string, separator string) []string {
	if s == "" {
		return []string{}
	}
	return splitTrim(s, separator)
}

// Helper function to split string by separator and trim spaces
func splitTrim(s string, separator string) []string {
	result := make([]string, 0)
	for _, item := range splitString(s, separator) {
		result = append(result, trim(item))
	}
	return result
}

// Helper function to trim spaces
func trim(s string) string {
	return trimSpace(s)
}

// Helper function to trim spaces
func trimSpace(s string) string {
	return trimAll(s, " ")
}

// Helper function to trim all occurrences of a character
func trimAll(s string, cutset string) string {
	if s == "" {
		return ""
	}
	return trimAllSpace(s, cutset)
}

// Helper function to trim all spaces
func trimAllSpace(s string, cutset string) string {
	if s == "" {
		return ""
	}
	return trimLeftRightSpace(s, cutset)
}

// Helper function to trim left and right spaces
func trimLeftRightSpace(s string, cutset string) string {
	if s == "" {
		return ""
	}
	return trimLeftSpace(trimRightSpace(s, cutset), cutset)
}

// Helper function to trim left spaces
func trimLeftSpace(s string, cutset string) string {
	if s == "" {
		return ""
	}
	return trimLeft(s, cutset)
}

// Helper function to trim right spaces
func trimRightSpace(s string, cutset string) string {
	if s == "" {
		return ""
	}
	return trimRight(s, cutset)
}

// Helper function to trim left occurrences of a character
func trimLeft(s string, cutset string) string {
	if s == "" {
		return ""
	}
	return trim(s)
}

// Helper function to trim right occurrences of a character
func trimRight(s string, cutset string) string {
	if s == "" {
		return ""
	}
	return s
}
