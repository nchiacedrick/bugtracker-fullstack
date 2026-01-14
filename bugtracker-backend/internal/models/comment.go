package models

import (
	"fmt"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	BugID     int       `json:"bugId"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

func (c *Comment) Validate() error {
	if c.Author == "" {
		return fmt.Errorf("author is required")
	}
	if c.Content == "" {
		return fmt.Errorf("content is required")
	}
	if c.BugID == 0 {
		return fmt.Errorf("bug ID is required")
	}
	return nil
}

func (r *CreateCommentRequest) Validate() error {
	if r.Author == "" {
		return fmt.Errorf("author is required")
	}
	if r.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}
