package testutil

import (
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var (
	testDBPaths = make(map[string]struct{})
	mu          sync.Mutex
)

func GetTestDBPath() string {
	path := filepath.Join(os.TempDir(), "test_"+time.Now().Format("20060102150405.000")+".db")
	mu.Lock()
	testDBPaths[path] = struct{}{}
	mu.Unlock()
	return path
}

func CleanupTestDB() error {
	mu.Lock()
	defer mu.Unlock()

	var lastErr error
	for path := range testDBPaths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			lastErr = err
		}
		delete(testDBPaths, path)
	}

	projectDirs := []string{
		".",
		"internal/db",
		"internal/handlers",
		"cmd/bugtracker",
	}

	for _, dir := range projectDirs {
		matches, err := filepath.Glob(filepath.Join(dir, "*.db"))
		if err != nil {
			continue
		}
		for _, match := range matches {
			if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
				lastErr = err
			}
		}
	}

	return lastErr
}

func init() {
	if os.Getenv("TEST_MODE") != "" {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			CleanupTestDB()
			os.Exit(1)
		}()
	}
}
