
package services

import (
	"fmt"
	"strings"
	"time"

	"database/sql"
	"Backend/internal/models"
)

type StockService struct {
	db *sql.DB
}

func NewStockService(db *sql.DB) *StockService {
	return &StockService{db: db}
}

func (s *StockService) GetStocks(filters models.StockFilters) (*models.StockResponse, error) {
	query := `
		SELECT id, ticker, company, brokerage, action, rating_from, rating_to, 
		       target_from, target_to, time, created_at, updated_at
		FROM stocks
	`
	
	args := []any{}
	argIndex := 1

	// Aplicar filtros
	if filters.Ticker != "" {
		query += fmt.Sprintf(" AND ticker ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Ticker+"%")
		argIndex++
	}

	if filters.Company != "" {
		query += fmt.Sprintf(" AND company ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Company+"%")
		argIndex++
	}

	if filters.Brokerage != "" {
		query += fmt.Sprintf(" AND brokerage ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Brokerage+"%")
		argIndex++
	}

	// Ordenamiento
	sortBy := "time"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	
	order := "DESC"
	if filters.Order == "asc" {
		order = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)

	// Paginación
	limit := 20
	if filters.Limit > 0 && filters.Limit <= 100 {
		limit = filters.Limit
	}

	offset := 0
	if filters.Page > 0 {
		offset = (filters.Page - 1) * limit
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(
			&stock.ID, &stock.Ticker, &stock.Company, &stock.Brokerage,
			&stock.Action, &stock.RatingFrom, &stock.RatingTo,
			&stock.TargetFrom, &stock.TargetTo, &stock.Time,
			&stock.CreatedAt, &stock.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return &models.StockResponse{Items: stocks}, nil
}

func (s *StockService) GetRecommendations() ([]models.Recommendation, error) {
	// Algoritmo simple de recomendación basado en:
	// 1. Upgrades recientes
	// 2. Targets altos
	// 3. Ratings positivos
	
	query := `
		SELECT ticker, company, rating_to, target_to, action, time
		FROM stocks
		WHERE rating_to IN ('Buy', 'Strong Buy', 'Outperform', 'Overweight')
		  AND action LIKE '%upgrade%'
		  AND time > NOW() - INTERVAL '30 days'
		GROUP BY ticker, company, rating_to, target_to, action, time
		ORDER BY time DESC
		LIMIT 10
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []models.Recommendation
	for rows.Next() {
		var ticker, company, rating, target, action string
		var time time.Time
		
		if err := rows.Scan(&ticker, &company, &rating, &target, &action, &time); err != nil {
			continue
		}

		score := calculateScore(rating, action, time)
		reason := generateReason(rating, action, target)

		recommendations = append(recommendations, models.Recommendation{
			Ticker:        ticker,
			Company:       company,
			Score:         score,
			Reason:        reason,
			TargetPrice:   target,
			CurrentRating: rating,
			Confidence:    score / 100,
		})
	}

	return recommendations, nil
}

func calculateScore(rating, action string, timestamp time.Time) float64 {
	score := 50.0 // Base score

	// Bonus por rating
	switch rating {
	case "Strong Buy":
		score += 30
	case "Buy":
		score += 20
	case "Outperform", "Overweight":
		score += 15
	}

	// Bonus por action
	if strings.Contains(strings.ToLower(action), "upgrade") {
		score += 20
	}

	// Bonus por recencia
	daysSince := time.Now().Sub(timestamp).Hours() / 24
	if daysSince < 7 {
		score += 10
	} else if daysSince < 14 {
		score += 5
	}

	if score > 100 {
		score = 100
	}

	return score
}

func generateReason(rating, action, target string) string {
	reasons := []string{}
	
	if strings.Contains(strings.ToLower(action), "upgrade") {
		reasons = append(reasons, "Reciente upgrade")
	}
	
	if rating == "Strong Buy" || rating == "Buy" {
		reasons = append(reasons, "Rating de compra")
	}
	
	if target != "" {
		reasons = append(reasons, "Precio objetivo: "+target)
	}
	
	if len(reasons) == 0 {
		return "Análisis técnico favorable"
	}
	
	return strings.Join(reasons, " • ")
}


// Implementación completa del método SyncAllData
func (s *StockService) SyncAllData(apiClient *APIClient) error {
	fmt.Println("Iniciando sincronización de datos...")
	
	stocks, err := apiClient.FetchAllStocks()
	if err != nil {
		return fmt.Errorf("error fetching stocks from API: %w", err)
	}

	fmt.Printf("Obtenidos %d registros de la API\n", len(stocks))

	if err := s.InsertStocks(stocks); err != nil {
		return fmt.Errorf("error inserting stocks into database: %w", err)
	}

	fmt.Printf("Sincronización completada: %d registros procesados\n", len(stocks))
	return nil
}

