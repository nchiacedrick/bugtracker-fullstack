package db

import (
	"bugtracker-backend/internal/models"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	cleanup := SetupTestDB(t)
	defer cleanup()

	bug := &models.Bug{
		Title:       "Test Bug",
		Description: "Test Description",
	}
	err := CreateBug(bug)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		bugID      string
		comment    *models.Comment
		shouldErr  bool
		errMessage string
	}{
		{
			name:  "Valid comment creation",
			bugID: strconv.Itoa(bug.ID),
			comment: &models.Comment{
				Author:  "Test User",
				Content: "Test Comment",
			},
			shouldErr: false,
		},
		{
			name:  "Invalid bug ID",
			bugID: "999",
			comment: &models.Comment{
				Author:  "Test User",
				Content: "Test Comment",
			},
			shouldErr:  true,
			errMessage: "bug not found",
		},
		{
			name:  "Invalid bug ID format",
			bugID: "invalid",
			comment: &models.Comment{
				Author:  "Test User",
				Content: "Test Comment",
			},
			shouldErr:  true,
			errMessage: "invalid bug ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateComment(tt.bugID, tt.comment)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.comment.ID)
				assert.NotEmpty(t, tt.comment.CreatedAt)
				expectedBugID, _ := strconv.Atoi(tt.bugID)
				assert.Equal(t, expectedBugID, tt.comment.BugID)
			}
		})
	}
}

func TestGetComments(t *testing.T) {
	cleanup := SetupTestDB(t)
	defer cleanup()

	bug := &models.Bug{
		Title:       "Test Bug",
		Description: "Test Description",
	}
	err := CreateBug(bug)
	assert.NoError(t, err)

	testComments := []*models.Comment{
		{Author: "User1", Content: "Comment 1"},
		{Author: "User2", Content: "Comment 2"},
	}

	for _, comment := range testComments {
		err := CreateComment(strconv.Itoa(bug.ID), comment)
		assert.NoError(t, err)
	}

	tests := []struct {
		name        string
		bugID       string
		wantErr     bool
		expectedErr string
	}{
		{
			name:    "Get existing comments",
			bugID:   strconv.Itoa(bug.ID),
			wantErr: false,
		},
		{
			name:        "Get comments for non-existent bug",
			bugID:       "999",
			wantErr:     true,
			expectedErr: "bug not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comments, err := GetComments(tt.bugID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, comments, len(testComments))
			}
		})
	}
}
