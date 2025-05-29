package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

// MigrateExplainTables creates tables for the explain module
func MigrateExplainTables(db *gorm.DB) error {
	// Create explain_requests table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS explain_requests (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			input_text TEXT NOT NULL,
			modes VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			
			CONSTRAINT fk_explain_user
				FOREIGN KEY(user_id)
				REFERENCES users(id)
				ON DELETE CASCADE
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create explain_requests table: %w", err)
	}
	
	// Create explain_results table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS explain_results (
			id SERIAL PRIMARY KEY,
			request_id INTEGER NOT NULL,
			mode VARCHAR(50) NOT NULL,
			result_content TEXT NOT NULL,
			audio_url VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			
			CONSTRAINT fk_result_request
				FOREIGN KEY(request_id)
				REFERENCES explain_requests(id)
				ON DELETE CASCADE
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create explain_results table: %w", err)
	}
	
	// Create auto_tags table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS auto_tags (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			category VARCHAR(50) NOT NULL,
			source INTEGER NOT NULL,
			count INTEGER DEFAULT 1,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			
			CONSTRAINT unique_tag_name_category UNIQUE(name, category)
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create auto_tags table: %w", err)
	}
	
	// Create explain_tags table
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS explain_tags (
			id SERIAL PRIMARY KEY,
			explain_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			
			CONSTRAINT fk_explain
				FOREIGN KEY(explain_id)
				REFERENCES explain_requests(id)
				ON DELETE CASCADE,
				
			CONSTRAINT fk_tag
				FOREIGN KEY(tag_id)
				REFERENCES auto_tags(id)
				ON DELETE CASCADE,
				
			CONSTRAINT unique_explain_tag UNIQUE(explain_id, tag_id)
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create explain_tags table: %w", err)
	}
	
	// Create indexes
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_explain_requests_user_id ON explain_requests(user_id)`).Error; err != nil {
		return fmt.Errorf("failed to create index on explain_requests: %w", err)
	}
	
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_explain_results_request_id ON explain_results(request_id)`).Error; err != nil {
		return fmt.Errorf("failed to create index on explain_results: %w", err)
	}
	
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_auto_tags_name ON auto_tags(name)`).Error; err != nil {
		return fmt.Errorf("failed to create index on auto_tags: %w", err)
	}
	
	return nil
}
