package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"bugtracker-backend/internal/models"

	_ "github.com/lib/pq"
)

var (
	db          *sql.DB
	initialized bool
)

func getDBURL() string {
	// Prefer explicit DATABASE_URL (one-line DSN). If not set, fall back to a sensible local default.
	if d := os.Getenv("DATABASE_URL"); d != "" {
		return d
	}
	return "postgres://postgres:postgres@localhost:5435/bugtracker?sslmode=disable"
}

func Init() error {
	if initialized {
		return fmt.Errorf("database already initialized")
	}

	dsn := getDBURL()
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if they do not exist
	schema := `
	CREATE TABLE IF NOT EXISTS bugs (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT,
		priority TEXT,
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ
	);
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		bug_id INTEGER REFERENCES bugs(id) ON DELETE CASCADE,
		content TEXT,
		author TEXT,
		created_at TIMESTAMPTZ
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	initialized = true
	return nil
}

func CreateBug(bug *models.Bug) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	if bug.CreatedAt.IsZero() {
		bug.CreatedAt = time.Now()
	}
	bug.UpdatedAt = time.Now()

	query := `INSERT INTO bugs (title, description, status, priority, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	err := db.QueryRow(query, bug.Title, bug.Description, bug.Status, bug.Priority, bug.CreatedAt, bug.UpdatedAt).Scan(&bug.ID)
	if err != nil {
		return fmt.Errorf("failed to insert bug: %w", err)
	}
	return nil
}

func GetBug(id int) (*models.Bug, error) {
	var bug models.Bug
	query := `SELECT id, title, description, status, priority, created_at, updated_at FROM bugs WHERE id=$1`
	row := db.QueryRow(query, id)
	if err := row.Scan(&bug.ID, &bug.Title, &bug.Description, &bug.Status, &bug.Priority, &bug.CreatedAt, &bug.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bug not found")
		}
		return nil, fmt.Errorf("failed to query bug: %w", err)
	}
	return &bug, nil
}

func GetAllBugs() ([]*models.Bug, error) {
	query := `SELECT id, title, description, status, priority, created_at, updated_at FROM bugs ORDER BY id`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query bugs: %w", err)
	}
	defer rows.Close()

	var bugs []*models.Bug
	for rows.Next() {
		var bug models.Bug
		if err := rows.Scan(&bug.ID, &bug.Title, &bug.Description, &bug.Status, &bug.Priority, &bug.CreatedAt, &bug.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan bug: %w", err)
		}
		bugs = append(bugs, &bug)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return bugs, nil
}

func DeleteBug(id int) error {
	res, err := db.Exec(`DELETE FROM bugs WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete bug: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("bug not found")
	}
	return nil
}

func Cleanup() {
	if db != nil {
		_ = db.Close()
		db = nil
	}
	initialized = false
}

func UpdateBug(bug *models.Bug) error {
	bug.UpdatedAt = time.Now()
	res, err := db.Exec(`UPDATE bugs SET title=$1, description=$2, status=$3, priority=$4, updated_at=$5 WHERE id=$6`, bug.Title, bug.Description, bug.Status, bug.Priority, bug.UpdatedAt, bug.ID)
	if err != nil {
		return fmt.Errorf("failed to update bug: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("bug not found")
	}
	return nil
}

func CleanupTestDB() error {
	if db == nil {
		return nil
	}
	// Truncate tables and restart identity for tests
	if _, err := db.Exec(`TRUNCATE TABLE comments, bugs RESTART IDENTITY CASCADE`); err != nil {
		return fmt.Errorf("failed to cleanup test db: %w", err)
	}
	return nil
}

func DeleteAllBugs() (int, error) {
	var count int
	row := db.QueryRow(`SELECT COUNT(*) FROM bugs`)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count bugs: %w", err)
	}
	if _, err := db.Exec(`DELETE FROM bugs`); err != nil {
		return 0, fmt.Errorf("failed to delete bugs: %w", err)
	}
	if _, err := db.Exec(`ALTER SEQUENCE bugs_id_seq RESTART WITH 1`); err != nil {
		// best-effort, ignore if sequence not present
	}
	return count, nil
}
