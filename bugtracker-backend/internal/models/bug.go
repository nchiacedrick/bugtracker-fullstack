package models

import (
	"fmt"
	"time"
)

type Bug struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateBugRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

func (b *Bug) Validate() error {
	if b.Title == "" {
		return fmt.Errorf("title is required")
	}
	if !isValidPriority(b.Priority) {
		return fmt.Errorf("invalid priority")
	}
	if !isValidStatus(b.Status) {
		return fmt.Errorf("invalid status")
	}
	return nil
}

func (r *CreateBugRequest) Validate() error {
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

func isValidPriority(p string) bool {
	validPriorities := []string{"Low", "Medium", "High"}
	return contains(validPriorities, p)
}

func isValidStatus(s string) bool {
	validStatuses := []string{"Open", "In Progress", "Closed"}
	return contains(validStatuses, s)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
