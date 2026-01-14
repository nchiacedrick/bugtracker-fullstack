package handlers

import (
	"bugtracker-backend/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response HealthResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, config.Backend_Version, response.Version)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
} 