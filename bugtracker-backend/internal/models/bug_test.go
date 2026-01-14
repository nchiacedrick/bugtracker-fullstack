package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBugValidation(t *testing.T) {
	tests := []struct {
		name    string
		bug     Bug
		isValid bool
		errMsg  string
	}{
		{
			name: "Valid bug",
			bug: Bug{
				Title:       "Test Bug",
				Description: "Test Description",
				Priority:    "High",
				Status:      "Open",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: true,
		},
		{
			name: "Missing title",
			bug: Bug{
				Description: "Test Description",
				Priority:    "High",
				Status:      "Open",
			},
			isValid: false,
			errMsg:  "title is required",
		},
		{
			name: "Invalid priority",
			bug: Bug{
				Title:       "Test Bug",
				Description: "Test Description",
				Priority:    "Invalid",
				Status:      "Open",
			},
			isValid: false,
			errMsg:  "invalid priority",
		},
		{
			name: "Invalid status",
			bug: Bug{
				Title:       "Test Bug",
				Description: "Test Description",
				Priority:    "High",
				Status:      "Invalid",
			},
			isValid: false,
			errMsg:  "invalid status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bug.Validate()
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCreateBugRequest(t *testing.T) {
	tests := []struct {
		name    string
		request CreateBugRequest
		isValid bool
		errMsg  string
	}{
		{
			name: "Valid request",
			request: CreateBugRequest{
				Title:       "Test Bug",
				Description: "Test Description",
				Priority:    "High",
				Status:      "Open",
			},
			isValid: true,
		},
		{
			name: "Missing title",
			request: CreateBugRequest{
				Description: "Test Description",
			},
			isValid: false,
			errMsg:  "title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}
