
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"Backend/internal/models"
)

type APIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type APIResponse struct {
	Items    []APIStock `json:"items"`
	NextPage string     `json:"next_page,omitempty"`
}

type APIStock struct {
	Ticker     string    `json:"ticker"`
	Company    string    `json:"company"`
	Brokerage  string    `json:"brokerage"`
	Action     string    `json:"action"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Time       time.Time `json:"time"`
}

func NewAPIClient(apiKey string, baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) FetchStocks(page string) (*APIResponse, error) {
	reqURL := c.baseURL
	
	if page != "" {
		u, err := url.Parse(reqURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing URL: %w", err)
		}
		
		q := u.Query()
		q.Set("next_page", page)
		u.RawQuery = q.Encode()
		reqURL = u.String()
	}

	fmt.Println("reqURL: ", reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("User-Agent", "karla/1.0")

	// fmt.Println("req: ", req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &apiResponse, nil
}

func (c *APIClient) FetchAllStocks() ([]models.Stock, error) {
	var allStocks []models.Stock
	nextPage := ""


for {
    response, err := c.FetchStocks(nextPage)
		// fmt.Println("response: ", response)
    if err != nil {
        return nil, fmt.Errorf("error fetching stocks: %w", err)
    }

    // Convert APIStock to models.Stock
    for _, apiStock := range response.Items {
        stock := models.Stock{
            Ticker:     apiStock.Ticker,
            Company:    apiStock.Company,
            Brokerage:  apiStock.Brokerage,
            Action:     apiStock.Action,
            RatingFrom: apiStock.RatingFrom,
            RatingTo:   apiStock.RatingTo,
            TargetFrom: apiStock.TargetFrom,
            TargetTo:   apiStock.TargetTo,
            Time:       apiStock.Time,
            CreatedAt:  time.Now(),
            UpdatedAt:  time.Now(),
        }

				// 
        allStocks = append(allStocks, stock)
    }

		  // Break the loop if there are no more pages
    if response.NextPage == "" {
        break
    }

    time.Sleep(500 * time.Millisecond) 
    nextPage = response.NextPage


    // Return results if we have accumulated a significant number of stocks
    if len(allStocks) >= 1000 {
        return allStocks, nil
    }
}


	return allStocks, nil
}


func (s *StockService) InsertStocks(stocks []models.Stock) error {

	if len(stocks) == 0 {
		return nil
	}

	query := `
		INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, 
		                   target_from, target_to, time, created_at, updated_at, score, reason, target_price, current_rating, confidence)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (ticker, brokerage) DO UPDATE SET
			action = EXCLUDED.action,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			created_at = EXCLUDED.created_at,
			updated_at = NOW(),
			time = EXCLUDED.time,
			score = EXCLUDED.score,
			reason = EXCLUDED.reason,
			target_price = EXCLUDED.target_price,
			current_rating = EXCLUDED.current_rating,
			confidence = EXCLUDED.confidence
	`

	tx, err := s.db.Begin()
	if err != nil {
		fmt.Printf("error starting transaction: %w", err)
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Printf("error preparing statement: %w", err)
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()


	for _, stock := range stocks {
		_, err := stmt.Exec(
			stock.Ticker, stock.Company, stock.Brokerage, stock.Action,
			stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
			stock.Time, stock.CreatedAt, stock.UpdatedAt, stock.Score, stock.Reason, stock.TargetPrice, stock.CurrentRating, stock.Confidence,
		)
		
		if err != nil {
			return fmt.Errorf("error inserting stock %s: %w", stock.Ticker, err)
		}
	}

	return tx.Commit()
}
