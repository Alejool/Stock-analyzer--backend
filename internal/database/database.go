
// internal/database/database.go
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
		UNIQUE(ticker, brokerage, time)
	);

	CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);
	CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks(company);
	CREATE INDEX IF NOT EXISTS idx_stocks_time ON stocks(time DESC);
	CREATE INDEX IF NOT EXISTS idx_stocks_rating_to ON stocks(rating_to);

	-- Insert sample data
	-- delete from stocks;
	TRUNCATE TABLE stocks;

	INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time)
	VALUES
		('AAPL', 'Apple Inc.', 'Morgan Stanley', 'upgrade', 'hold', 'buy', '150', '180', NOW()),
		('GOOGL', 'Alphabet Inc.', 'Goldman Sachs', 'downgrade', 'buy', 'hold', '2800', '2600', NOW()),
		('MSFT', 'Microsoft Corp.', 'JP Morgan', 'reiterate', 'buy', 'buy', '310', '320', NOW()),
		('AMZN', 'Amazon.com Inc.', 'Citigroup', 'initiate', 'sell', 'buy', 'null', '3500', NOW()),
		('TSLA', 'Tesla Inc.', 'Bank of America', 'upgrade', 'sell', 'hold', '200', '250', NOW())
	ON CONFLICT (ticker, brokerage, time) DO NOTHING;
	`
	_, err := db.Exec(query)
	return err
}