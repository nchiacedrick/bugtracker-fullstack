package db

import (
	"bugtracker-backend/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func CreateComment(bugID string, comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	var err error
	comment.BugID, err = strconv.Atoi(bugID)
	if err != nil {
		return fmt.Errorf("invalid bug ID format")
	}

	// ensure bug exists
	if _, err := GetBug(comment.BugID); err != nil {
		return fmt.Errorf("bug not found")
	}

	query := `INSERT INTO comments (bug_id, content, author, created_at) VALUES ($1,$2,$3,$4) RETURNING id`
	if err := db.QueryRow(query, comment.BugID, comment.Content, comment.Author, comment.CreatedAt).Scan(&comment.ID); err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}
	return nil
}

func GetComments(bugID string) ([]models.Comment, error) {
	var comments []models.Comment
	bugIDInt, err := strconv.Atoi(bugID)
	if err != nil {
		return nil, fmt.Errorf("invalid bug ID format")
	}

	if _, err := GetBug(bugIDInt); err != nil {
		return nil, fmt.Errorf("bug not found")
	}

	rows, err := db.Query(`SELECT id, bug_id, content, author, created_at FROM comments WHERE bug_id=$1 ORDER BY id`, bugIDInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return comments, nil
		}
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.BugID, &c.Content, &c.Author, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return comments, nil
}
