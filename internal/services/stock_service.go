package services

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"Backend/internal/models"
	"database/sql"
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
		       target_from, target_to, time, created_at, updated_at, score, confidence,
		       count(*) OVER() AS total_register,
		       count(CASE WHEN rating_to = 'Buy' THEN 1 END) OVER() AS buy_count,
		       (SELECT COUNT(DISTINCT brokerage) FROM stocks) AS total_brokerages,
		       max(updated_at) OVER() AS last_update
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

	if filters.ProductID != 0 {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, filters.ProductID)
		argIndex++
	}

	if filters.Score > 0 {
		query += fmt.Sprintf(" AND score >= $%d", argIndex)
		args = append(args, filters.Score)
		argIndex++
	}

if filters.Confidence != "" {
    if strings.ToUpper(filters.Confidence) == "ASC" {
        query += " ORDER BY confidence ASC"
    } else if strings.ToUpper(filters.Confidence) == "DESC" {
        query += " ORDER BY confidence DESC"
    }
}

	// fmt.Println("query: ", query)

	// Ordenamiento
	sortBy := "confidence"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}

	order := ""
	if filters.Order == "asc" {
		order = "ASC"
	}
	if filters.Order == "desc" {
		order = "DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)

	// Paginación


	
	limit := 0
	if filters.Limit > 0 {
		limit = filters.Limit
	}

	// offset := 0
	// if filters.Page > 0 {
	// 	offset = (filters.Page - 1) * limit
	// }

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d ", limit)
	}

	
	

	// fmt.Println("query: ", query)

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
			&stock.CreatedAt, &stock.UpdatedAt, &stock.Score, &stock.Confidence,
			&stock.TotalRegister, &stock.BuyCount, &stock.TotalBrokerages, &stock.LastUpdateFilter,
		)
		if err != nil {
			return nil, err
		}

		stocks = append(stocks, stock)
	}

	// oragnizata stocks por confindece
	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].Confidence > stocks[j].Confidence
	})

	return &models.StockResponse{Items: stocks}, nil
}

func (s *StockService) GetRecommendations() ([]models.Stock, error) {
	// Algoritmo simple de recomendación basado en:
	// 1. Upgrades recientes
	// 2. Targets altos
	// 3. Ratings positivos

	query := `
		SELECT ticker, company, rating_to, brokerage, target_to,  rating_from, action, time, score, confidence
		FROM stocks
		WHERE 
		  time > NOW() - INTERVAL '30 days'
		  AND score > 0
		ORDER BY confidence DESC, score DESC
		LIMIT 1
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []models.Stock
	for rows.Next() {
		var recomendation models.Stock

		err := rows.Scan(
			&recomendation.Ticker, 
			&recomendation.Company,
			&recomendation.RatingTo, 
			&recomendation.Brokerage, 
			&recomendation.TargetTo, 
			&recomendation.RatingFrom,
			&recomendation.Action, 
			 &recomendation.Time,
			&recomendation.Score, 
			&recomendation.Confidence,
		)
		if err != nil {
			return nil, err
		}
		
		recommendations = append(recommendations, recomendation)
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

	// fmt.Println("stocks: ", stocks)
	if err != nil {
		return fmt.Errorf("error fetching stocks from API: %w", err)
	}

	// agregar score y confidence
	for i := range stocks {
		score := calculateScore(stocks[i].RatingTo, stocks[i].Action, stocks[i].Time)
		reason := generateReason(stocks[i].RatingTo, stocks[i].Action, stocks[i].TargetTo)

		stocks[i].Score = score
		stocks[i].Reason = reason
		stocks[i].CurrentRating = stocks[i].RatingTo
		stocks[i].Confidence = score / 100
	}

	fmt.Print("score y confidence agregados")

	// fmt.Printf("Obtenidos %d registros de la API\n", len(stocks))

	if err := s.InsertStocks(stocks); err != nil {
		return fmt.Errorf("error inserting stocks into database: %w", err)
	}
	

	fmt.Printf("Sincronización completada: %d registros procesados\n", len(stocks))
	return nil
}
