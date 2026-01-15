package db

import (
	"os"
	"testing"

	"bugtracker-backend/internal/models"
	"bugtracker-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseInitialization(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping DB tests")
	}
	os.Setenv("DATABASE_URL", url)
	defer os.Unsetenv("DATABASE_URL")

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer func() {
		_ = CleanupTestDB()
		Cleanup()
	}()

	bug := &models.Bug{Title: "Test", Description: "Test"}
	err := CreateBug(bug)
	assert.NoError(t, err)
}

func TestMultipleInitializations(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping DB tests")
	}
	os.Setenv("DATABASE_URL", url)
	defer os.Unsetenv("DATABASE_URL")

	err := Init()
	assert.NoError(t, err)
	Cleanup()

	err = Init()
	assert.NoError(t, err)
	defer func() {
		_ = CleanupTestDB()
		Cleanup()
	}()
}

func TestCleanup(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping DB tests")
	}
	os.Setenv("DATABASE_URL", url)
	defer os.Unsetenv("DATABASE_URL")

	err := Init()
	assert.NoError(t, err)
	Cleanup()

	// Test DB is inaccessible after cleanup
	bug := &models.Bug{Title: "Test", Description: "Test"}
	err = CreateBug(bug)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database not initialized",
		"Should get 'database not initialized' error after cleanup")
}

func TestInitWithInvalidPath(t *testing.T) {
	// Setting an invalid DATABASE_URL should cause Init to fail
	original := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", original)

	invalid := "postgres://invalid:invalid@localhost:5432/nope?sslmode=disable"
	os.Setenv("DATABASE_URL", invalid)

	err := Init()
	assert.Error(t, err)
	// Accept either open or ping failures from the driver (auth, network, etc.)
	assert.Contains(t, err.Error(), "failed to")
	defer func() {
		_ = CleanupTestDB()
		Cleanup()
	}()
}

func TestConcurrentInitializations(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping DB tests")
	}
	os.Setenv("DATABASE_URL", url)
	defer os.Unsetenv("DATABASE_URL")

	err := Init()
	assert.NoError(t, err)
	defer func() {
		_ = CleanupTestDB()
		Cleanup()
	}()

	err = Init()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database already initialized")
}

func TestMain(m *testing.M) {
	os.Setenv("TEST_MODE", "1")

	code := m.Run()

	testutil.CleanupTestDB()

	os.Exit(code)
}
