package handlers

import (
	"bugtracker-backend/internal/db"
	"testing"
)

func setupTestDB(t *testing.T) func() {
	cleanup := db.SetupTestDB(t)
	return cleanup
} 