package db

import (
	"os"
	"testing"
)

func SetupTestDB(t *testing.T) func() {
	// Tests rely on a database being initialized. Ensure Init runs and tables are clean.
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	if err := CleanupTestDB(); err != nil {
		t.Fatalf("Failed to cleanup test database: %v", err)
	}

	return func() {
		Cleanup()
		// no file to remove for Postgres
		_ = os.Getenv("TEST_DATABASE_URL")
	}
}
