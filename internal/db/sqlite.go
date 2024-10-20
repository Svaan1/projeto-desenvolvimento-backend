package db

import (
	"context"
	"database/sql"
	"encoding/json"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	Client *sql.DB
}

func NewSQLiteDB(ctx context.Context, dbFile string) (*SQLiteDB, error) {
	// Open a SQLite database file
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// Optionally, create a table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS objects (
		key TEXT PRIMARY KEY,
		value TEXT
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return &SQLiteDB{Client: db}, nil
}

// SetObject stores an object as a JSON string in SQLite
func (s *SQLiteDB) SetObject(ctx context.Context, key string, obj interface{}) error {
	// Marshal the object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err // JSON marshaling error
	}

	// Upsert: insert the object or update if it already exists
	_, err = s.Client.ExecContext(ctx, `INSERT INTO objects (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, data)
	if err != nil {
		return err // SQL error
	}

	return nil
}

// GetObject retrieves a whole object from SQLite by its key
func (s *SQLiteDB) GetObject(ctx context.Context, key string, obj interface{}) error {
	// Get the JSON string from SQLite
	var value string
	err := s.Client.QueryRowContext(ctx, `SELECT value FROM objects WHERE key = ?`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // Key does not exist
		}
		return err // Other error
	}

	// Unmarshal the JSON string into the provided obj
	err = json.Unmarshal([]byte(value), obj)
	if err != nil {
		return err // Unmarshal error
	}

	return nil // Successfully retrieved
}

// Close closes the SQLite database connection
func (s *SQLiteDB) Close() error {
	return s.Client.Close()
}
