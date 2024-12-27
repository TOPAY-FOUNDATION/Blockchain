package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// InitDB initializes a SQLite database for the blockchain
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create table for blockchain storage
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS blocks (
		index INTEGER PRIMARY KEY,
		data TEXT
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	fmt.Printf("Database initialized at %s\n", dbPath)
	return db, nil
}

// SaveBlockToDB saves a block to the database
func SaveBlockToDB(db *sql.DB, index int, block interface{}) error {
	jsonData, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %v", err)
	}

	_, err = db.Exec("INSERT OR REPLACE INTO blocks (index, data) VALUES (?, ?)", index, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to save block to database: %v", err)
	}

	fmt.Printf("Block %d saved to database\n", index)
	return nil
}

// LoadBlocksFromDB loads all blocks from the database
func LoadBlocksFromDB(db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT data FROM blocks ORDER BY index")
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %v", err)
	}
	defer rows.Close()

	var blocks []map[string]interface{}
	for rows.Next() {
		var jsonData string
		if err := rows.Scan(&jsonData); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		var block map[string]interface{}
		if err := json.Unmarshal([]byte(jsonData), &block); err != nil {
			return nil, fmt.Errorf("failed to unmarshal block: %v", err)
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}
