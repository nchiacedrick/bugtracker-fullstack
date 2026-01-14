package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommentValidation(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
		isValid bool
		errMsg  string
	}{
		{
			name: "Valid comment",
			comment: Comment{
				ID:       1,
				BugID:    1,
				Author:   "Test Author",
				Content:  "Test Content",
				CreatedAt: time.Now(),
			},
			isValid: true,
		},
		{
			name: "Missing author",
			comment: Comment{
				BugID:  1,
				Content: "Test Content",
			},
			isValid: false,
			errMsg:  "author is required",
		},
		{
			name: "Missing content",
			comment: Comment{
				BugID:  1,
				Author: "Test Author",
			},
			isValid: false,
			errMsg:  "content is required",
		},
		{
			name: "Missing bug ID",
			comment: Comment{
				Author:  "Test Author",
				Content: "Test Content",
			},
			isValid: false,
			errMsg:  "bug ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCreateCommentRequest(t *testing.T) {
	tests := []struct {
		name    string
		request CreateCommentRequest
		isValid bool
		errMsg  string
	}{
		{
			name: "Valid request",
			request: CreateCommentRequest{
				Author:  "Test Author",
				Content: "Test Content",
			},
			isValid: true,
		},
		{
			name: "Missing author",
			request: CreateCommentRequest{
				Content: "Test Content",
			},
			isValid: false,
			errMsg:  "author is required",
		},
		{
			name: "Missing content",
			request: CreateCommentRequest{
				Author: "Test Author",
			},
			isValid: false,
			errMsg:  "content is required",
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
