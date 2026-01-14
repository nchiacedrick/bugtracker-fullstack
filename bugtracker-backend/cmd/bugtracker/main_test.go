package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"bugtracker-backend/internal/handlers"
	"bugtracker-backend/internal/testutil"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("TEST_MODE", "1")
	code := m.Run()
	testutil.CleanupTestDB()
	os.Exit(code)
}

func TestServerInitialization(t *testing.T) {
	os.Setenv("DB_PATH", testutil.GetTestDBPath())
	defer testutil.CleanupTestDB()

	testPort := ":8081" // Use a different port for testing
	srv := createTestServer()
	srv.Addr = testPort // Override the server port

	// Start the server in a goroutine
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Errorf("Expected ErrServerClosed, got %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test CORS with an OPTIONS request
	req, err := http.NewRequest("OPTIONS", "http://localhost:8081/api/health", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))

	// Test graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestCORSConfiguration(t *testing.T) {
	os.Setenv("DB_PATH", testutil.GetTestDBPath())
	defer testutil.CleanupTestDB()

	testPort := ":8081"
	srv := createTestServer()
	srv.Addr = testPort

	// Start the server
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Errorf("Expected ErrServerClosed, got %v", err)
		}
	}()
	defer func() {
		srv.Shutdown(context.Background())
	}()

	time.Sleep(100 * time.Millisecond)

	tests := []struct {
		name              string
		method            string
		origin            string
		expectedStatus    int
		expectedOrigin    string
		requestHeaders    map[string]string
		expectedHeaders   map[string]string
		shouldHaveHeaders bool
	}{
		{
			name:           "Valid OPTIONS request from allowed origin",
			method:         "OPTIONS",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusNoContent,
			expectedOrigin: "http://localhost:3000",
			requestHeaders: map[string]string{
				"Access-Control-Request-Method": "GET",
			},
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Methods": "GET",
			},
			shouldHaveHeaders: true,
		},
		{
			name:           "OPTIONS request from disallowed origin",
			method:         "OPTIONS",
			origin:         "http://evil.com",
			expectedStatus: http.StatusNoContent,
			expectedOrigin: "",
			requestHeaders: map[string]string{
				"Access-Control-Request-Method": "GET",
			},
			shouldHaveHeaders: false,
		},
		{
			name:              "Valid GET request from allowed origin",
			method:            "GET",
			origin:            "http://localhost:3000",
			expectedStatus:    http.StatusOK,
			expectedOrigin:    "http://localhost:3000",
			requestHeaders:    map[string]string{},
			shouldHaveHeaders: true,
		},
		{
			name:           "OPTIONS request with invalid method",
			method:         "OPTIONS",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusNoContent,
			expectedOrigin: "",
			requestHeaders: map[string]string{
				"Access-Control-Request-Method": "INVALID",
			},
			shouldHaveHeaders: false,
		},
		{
			name:           "Request with credentials",
			method:         "GET",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusOK,
			expectedOrigin: "http://localhost:3000",
			requestHeaders: map[string]string{
				"Cookie": "session=123",
			},
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Credentials": "true",
			},
			shouldHaveHeaders: true,
		},
	}

	client := &http.Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, fmt.Sprintf("http://localhost%s/api/health", testPort), nil)
			assert.NoError(t, err)

			req.Header.Set("Origin", tt.origin)
			for k, v := range tt.requestHeaders {
				req.Header.Set(k, v)
			}

			// Log request details
			t.Logf("Request: %s %s", tt.method, req.URL)
			t.Logf("Request Headers: %+v", req.Header)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Log response details
			t.Logf("Response Status: %d", resp.StatusCode)
			t.Logf("Response Headers: %+v", resp.Header)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedOrigin, resp.Header.Get("Access-Control-Allow-Origin"))

			if tt.shouldHaveHeaders {
				// Log each header we're checking
				t.Logf("Checking CORS headers:")
				t.Logf("Access-Control-Allow-Origin: %s", resp.Header.Get("Access-Control-Allow-Origin"))
				t.Logf("Access-Control-Allow-Methods: %s", resp.Header.Get("Access-Control-Allow-Methods"))
				t.Logf("Access-Control-Allow-Headers: %s", resp.Header.Get("Access-Control-Allow-Headers"))

				// Check for presence of basic CORS headers
				assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"))

				// Check expected headers if specified
				for k, v := range tt.expectedHeaders {
					assert.Equal(t, v, resp.Header.Get(k))
				}

				// For preflight requests, verify CORS preflight headers
				if tt.method == "OPTIONS" {
					assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Methods"))
					assert.NotEmpty(t, resp.Header.Get("Access-Control-Max-Age"))
				}
			}
		})
	}
}

// Test-specific server creation
func createTestServer() *http.Server {
	r := mux.NewRouter()
	// Health check endpoint before CORS
	r.HandleFunc("/api/health", handlers.HealthCheck).Methods("GET")

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true,
	}).Handler(r)

	return &http.Server{
		Addr:    ":8081",
		Handler: handler,
	}
}
