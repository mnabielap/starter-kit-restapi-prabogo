package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUserToken, downUserToken)
}

func upUserToken(ctx context.Context, tx *sql.Tx) error {
	// Users Table
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(50) DEFAULT 'user',
		is_email_verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);`)
	if err != nil {
		return err
	}

	// Tokens Table
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS tokens (
		id SERIAL PRIMARY KEY,
		token VARCHAR(255) NOT NULL,
		user_id VARCHAR(36) NOT NULL,
		type VARCHAR(50) NOT NULL,
		expires TIMESTAMP NOT NULL,
		blacklisted BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return err
	}

	// Indexes for performance
	_, err = tx.Exec(`CREATE INDEX idx_tokens_token ON tokens(token);`)
	if err != nil {
		return err
	}

	return nil
}

func downUserToken(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE tokens;`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DROP TABLE users;`)
	if err != nil {
		return err
	}
	return nil
}