package db

import (
	"bugtracker-backend/internal/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

func CreateComment(bugID string, comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	comment.ID = int(uuid.New().ID())
	var err error
	comment.BugID, err = strconv.Atoi(bugID)
	if err != nil {
		return fmt.Errorf("invalid bug ID format")
	}

	_, err = GetBug(comment.BugID)
	if err != nil {
		return fmt.Errorf("bug not found")
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(commentsBucket)
		encoded, err := json.Marshal(comment)
		if err != nil {
			return fmt.Errorf("failed to marshal comment: %v", err)
		}
		return b.Put(itob(comment.ID), encoded)
	})
}

func GetComments(bugID string) ([]models.Comment, error) {
	var comments []models.Comment
	bugIDInt, err := strconv.Atoi(bugID)
	if err != nil {
		return nil, fmt.Errorf("invalid bug ID format")
	}

	_, err = GetBug(bugIDInt)
	if err != nil {
		return nil, fmt.Errorf("bug not found")
	}

	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(commentsBucket)
		return b.ForEach(func(k, v []byte) error {
			var comment models.Comment
			if err := json.Unmarshal(v, &comment); err != nil {
				return err
			}
			if comment.BugID == bugIDInt {
				comments = append(comments, comment)
			}
			return nil
		})
	})

	return comments, err
}
