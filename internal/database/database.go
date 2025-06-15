
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexión: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	query := `
		-- DROP TABLE IF EXISTS stocks;

	CREATE TABLE IF NOT EXISTS stocks (
		id SERIAL PRIMARY KEY,
		ticker VARCHAR(10) NOT NULL,
		company VARCHAR(255) NOT NULL,
		brokerage VARCHAR(255) NOT NULL,
		action VARCHAR(50) NOT NULL,
		rating_from VARCHAR(50),
		rating_to VARCHAR(50),
		target_from VARCHAR(20),
		target_to VARCHAR(20),
		time TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		score FLOAT,
		reason VARCHAR(255),
		target_price VARCHAR(20),
		current_rating VARCHAR(50),
		confidence FLOAT,
		UNIQUE(ticker, brokerage, time)
	);

	CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);
	CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks(company);
	CREATE INDEX IF NOT EXISTS idx_stocks_time ON stocks(time DESC);
	CREATE INDEX IF NOT EXISTS idx_stocks_rating_to ON stocks(rating_to);

	-- delete from stocks;
	-- TRUNCATE TABLE stocks;

	 

	`
	_, err := db.Exec(query)
	return err
}