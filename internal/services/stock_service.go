package services

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"strconv"

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


	fmt.Println("filters: ", filters)

	args := []any{}
	argIndex := 1

	query += fmt.Sprintf(" WHERE 1=1 ")

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

		if filters.Today== "true" {
		query += fmt.Sprintf(" AND DATE(time) >= CURRENT_DATE")
		argIndex++
	}

	if filters.Confidence != "" {
		if strings.ToUpper(filters.Confidence) == "ASC" {
			query += " ORDER BY confidence ASC"
		} else if strings.ToUpper(filters.Confidence) == "DESC" {
			query += " ORDER BY confidence DESC"
		}
	}


	// Ordenamiento
	sortBy := "confidence"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}

	order := "DESC"
	if filters.Order == "ASC" {
		order = "ASC"
	}
	if filters.Order == "DESC" {
		order = "DESC"
	}

	query += fmt.Sprintf(" GROUP BY id, ticker, brokerage ORDER BY %s %s  ", sortBy, order)

	// fmt.Println("query: ", query)

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


		fmt.Println("query: ", query)

	// Log the SQL query and parameters for debugging API requests
	// fmt.Printf("SQL Query: %s\nParameters: %v\n", query, args)

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

	// organizate data stocks por confindece
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
		WITH today_records AS (
			SELECT *
			FROM stocks 
			WHERE DATE(time) = CURRENT_DATE
			AND score > 0
		),
		yesterday_records AS (
			SELECT *
			FROM stocks
			WHERE DATE(time) = CURRENT_DATE - INTERVAL '1 day'
			AND score > 0
		)
		SELECT ticker, company, rating_to, brokerage, target_to, rating_from, action, time, score, confidence
		FROM (
			SELECT *
			FROM today_records
			UNION ALL
			SELECT *
			FROM yesterday_records
			WHERE NOT EXISTS (SELECT 1 FROM today_records)
		) combined_records
		ORDER BY confidence DESC, score DESC, time DESC
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
func calculateScore(ratingFrom, ratingTo, action, targetFromStr, targetToStr string, timestamp time.Time) float64 {
	score := 50.0

	// Mapeo completo de ratings con sus rankings
	ratingRank := map[string]int{
		// Ratings negativos
		"sell":                    1,
		"underperform":           2,
		"sector underperform":    3,
		"underweight":            4,
		
		// Ratings neutrales
		"hold":                   5,
		"neutral":                5,
		"equal weight":           5,
		"market perform":         5,
		"sector perform":         5,
		"in-line":                5,
		"peer perform":           5,
		"sector weight":          5,
		
		// Ratings positivos bajos
		"positive":               6,
		"outperformer":           6,
		
		// Ratings positivos medios
		"outperform":             7,
		"market outperform":      7,
		"sector outperform":      7,
		
		// Ratings positivos altos
		"overweight":             8,
		"buy":                    8,
		
		// Ratings más altos
		"strong-buy":             9,
		"speculative buy":        9,
	}

	// Procesar ratings
	rf := strings.ToLower(strings.TrimSpace(ratingFrom))
	rt := strings.ToLower(strings.TrimSpace(ratingTo))

	if fromRank, ok1 := ratingRank[rf]; ok1 {
		if toRank, ok2 := ratingRank[rt]; ok2 {
			delta := toRank - fromRank
			score += float64(delta) * 3
		}
	}

	// Procesar targets
	targetFromStr = strings.ReplaceAll(strings.ReplaceAll(targetFromStr, "$", ""), ",", "")
	targetToStr = strings.ReplaceAll(strings.ReplaceAll(targetToStr, "$", ""), ",", "")
	targetFromFloat, err1 := strconv.ParseFloat(targetFromStr, 64)
	targetToFloat, err2 := strconv.ParseFloat(targetToStr, 64)
	targetFrom := float64(int64(targetFromFloat*100)) / 100
	targetTo := float64(int64(targetToFloat*100)) / 100

	if err1 == nil && err2 == nil && targetFrom > 0 {
		percentChange := (targetTo - targetFrom) / targetFrom
		score += percentChange * 100 // +10% => +10 puntos
	}

	actionLower := strings.ToLower(strings.TrimSpace(action))
	
	actionLower = strings.TrimSuffix(actionLower, " by")
	
	switch actionLower {
	case "upgraded":
		score += 15
	case "downgraded":
		score -= 15
	case "initiated":
		score += 8
	case "target raised":
		score += 8
	case "target lowered":
		score -= 8
	case "reiterated":
		score += 2  
	case "target set":
		score += 5  
	}

daysSince := time.Since(timestamp).Hours() / 24
if daysSince < 1 {
    score += 10          
} else if daysSince < 2 {
    score += 2            
} else if daysSince < 3 {
    score -= 2           
} else if daysSince < 5 {
    score -= 15         
} else if daysSince < 8 {
    score = 40          
} else {
    score = 30          
}

	// Limitar el score entre 0 y 100
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
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

func (s *StockService) SyncAllData(apiClient *APIClient) error {
	fmt.Println("Iniciando sincronización de datos...")

	stocks, err := apiClient.FetchAllStocks()

	// fmt.Println("stocks: ", stocks)
	if err != nil {
		return fmt.Errorf("error fetching stocks from API: %w", err)
	}

	// agregar score y confidence
	for i := range stocks {
		score := calculateScore(stocks[i].RatingFrom, stocks[i].RatingTo, stocks[i].Action, stocks[i].TargetFrom, stocks[i].TargetTo, stocks[i].Time)
		reason := generateReason(stocks[i].RatingTo, stocks[i].Action, stocks[i].TargetTo)

		stocks[i].Score = float64(int64(score*100)) / 100
		stocks[i].Reason = reason
		stocks[i].CurrentRating = stocks[i].RatingTo
		stocks[i].Confidence = float64(int64((score/100)*1000)) / 1000
	}

	fmt.Print("score y confidence agregados")

	// fmt.Printf("Obtenidos %d registros de la API\n", len(stocks))

	if err := s.InsertStocks(stocks); err != nil {
		return fmt.Errorf("error inserting stocks into database: %w", err)
	}

	fmt.Printf("Sincronización completada: %d registros procesados\n", len(stocks))
	return nil
}
