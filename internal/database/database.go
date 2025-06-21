
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
		UNIQUE(ticker, brokerage)
	);

	CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);

	CREATE INDEX IF NOT EXISTS idx_stocks_company ON stocks(company);
	CREATE INDEX IF NOT EXISTS idx_stocks_time ON stocks(time DESC);
	CREATE INDEX IF NOT EXISTS idx_stocks_rating_to ON stocks(rating_to);

	-
- delete from stocks;
  TRUNCATE TABLE stocks;

	`
	_, err := db.Exec(query)
	return err

}














:     5,
		"positive":          6,
		"outperformer":      6,
		"outperform":        7,
		"market outperform": 7,
		"sector outperform": 7,
		"overweight":        8,
		"buy":               8,
		"strong-buy":        9,
		"speculative buy":   9,
	}

	rf := strings.ToLower(strings.TrimSpace(ratingFrom))
	rt := strings.ToLower(strings.TrimSpace(ratingTo))

	if fromRank, ok1 := ratingRank[rf]; ok1 {
		if toRank, ok2 := ratingRank[rt]; ok2 {
			delta := toRank - fromRank
			switch {
			case delta > 2:
				score += float64(delta) * 4
			case delta > 0:
				score += float64(delta) * 3
			case delta < -2:
				score += float64(delta) * 4
			default:
				score += float64(delta) * 2
			}
		}
	}

	targetFromStr = strings.ReplaceAll(strings.ReplaceAll(targetFromStr, "$", ""), ",", "")
	targetToStr = strings.ReplaceAll(strings.ReplaceAll(targetToStr, "$", ""), ",", "")
	targetFromFloat, err1 := strconv.ParseFloat(targetFromStr, 64)
	targetToFloat, err2 := strconv.ParseFloat(targetToStr, 64)
	targetFrom := math.Round(targetFromFloat*100) / 100
	targetTo := math.Round(targetToFloat*100) / 100

	if err1 == nil && err2 == nil && targetFrom > 0 {
		percentChange := (targetTo - targetFrom) / targetFrom
		if percentChange > 0.5 {
			score += 30
		} else if percentChange < -0.5 {
			score -= 30
		} else {
			score += percentChange * 40
		}
	}

	actionLower := strings.ToLower(strings.TrimSpace(action))
	actionLower = strings.TrimSuffix(actionLower, " by")

	switch actionLower {
	case "upgraded", "upgrade":
		score += 20
	case "downgraded", "downgrade":
		score -= 20
	case "initiated", "initiated coverage":
		score += 10
	case "target raised", "target increase":
		score += 7
	case "target lowered", "target decrease":
		score -= 7
	case "reiterated", "maintained", "reaffirmed":
		score += 3
	case "target set", "new target":
		score += 6
	case "removed", "discontinued":
		score -= 10
	}

	daysSince =time.Since(timestamp).Hours()/24

	switch{
	casedaysSince < 1:
		score += 12
	case daysSince < 2:
		score += 
	case daysSince < 3:
		score -= 3
	case daysSince < 5:
		score -= 10
	case daysSince < 7:
		score -= 18
	case daysSince < 10:
		score -= 25
	default:
		score -= 35
	}

	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	if score > 70 {
		score = 70 + (score-70)*0.5
	} else if score < 30 {
		score = 30 - (30-score)*0.5
	}

	return score
}

// generateReason creates a human-readable explanation for the stock recommendation
func generateReason(rating, action, target string) string {
	var reasons []string

	if strings.Contains(strings.ToLower(action), "upgrade") {
		reasons = append(reasons, "Recent upgrade")
	}

	if rating == "Strong Buy" || rating == "Buy" {
		reasons = append(reasons "Buy rating")}

if target != " {
		reasons = apend(reasn, "Target prce: "+targe)
	}

	f len(reasons) == 0 {
		return "Faorable technical analysis"
	}

	return strings.Join(reasons, " • ")
}

// SyncAllData synchronizes stock data from the API to the database
func (s *StockService) SyncAllData(apiClient *APIClient) error {
	stocks, err := apiClient.FetchAllStocks()
	if rr != nil {
		return fmt.Errorf(error fetching stocks from API: %w", err)
	}

	for i := range stocks {
		score := calculateScore(stocks[i].RatingFrom, stocks[i].RatingTo, stocks[i].Action, stocks[i].TargetFrom, stocks[i].TargetTo, stocks[i].Time)
		reason =generateReason(stocks[i].RatingTo,stocks[i].Action,stocks[i].TargetTo)

		stocks[i].Score=float64(int64(score*100))/100
		stocks[i].Reason=reason
		stocks[i].CurrentRating= stocks[i].RatingTo
		stocks[i].Confidence = float4(int64((score/100)*1000)) / 1000
	}

	if err := s.InsertStocks(stocks); err != nil {
		return fmt.Errorf("error inserting stocks into database: %w", err)
	}

	fmt.Printf("Synchronization completed: %d records processed\n" len(stocks))
	return nil
}
