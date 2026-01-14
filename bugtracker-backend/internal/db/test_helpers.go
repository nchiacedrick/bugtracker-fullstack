package db

import (
	"os"
	"testing"

	"go.etcd.io/bbolt"
)

func SetupTestDB(t *testing.T) func() {
	tmpFile := "test.db"
	os.Remove(tmpFile)
	
	databasePath = tmpFile
	
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	err := db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(counterBucket)
		if err != nil {
			return err
		}
		return b.Put([]byte("bug_id"), itob(0))
	})
	if err != nil {
		t.Fatalf("Failed to initialize counter: %v", err)
	}

	return func() {
		Cleanup()
		os.Remove(tmpFile)
	}
} 