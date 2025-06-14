
// internal/models/stock.go
package models

import (
	"time"
)

type Stock struct {
	ID         int       `json:"id" db:"id"`
	Ticker     string    `json:"ticker" db:"ticker"`
	Company    string    `json:"company" db:"company"`
	Brokerage  string    `json:"brokerage" db:"brokerage"`
	Action     string    `json:"action" db:"action"`
	RatingFrom string    `json:"rating_from" db:"rating_from"`
	RatingTo   string    `json:"rating_to" db:"rating_to"`
	TargetFrom string    `json:"target_from" db:"target_from"`
	TargetTo   string    `json:"target_to" db:"target_to"`
	Time       time.Time `json:"time" db:"time"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StockResponse struct {
	Items    []Stock `json:"items"`
	NextPage string  `json:"next_page,omitempty"`
}

type StockFilters struct {
	Ticker    string `json:"ticker" form:"ticker"`
	Company   string `json:"company" form:"company"`
	Brokerage string `json:"brokerage" form:"brokerage"`
	Action    string `json:"action" form:"action"`
	Rating    string `json:"rating" form:"rating"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	Order     string `json:"order" form:"order"`
	Page      int    `json:"page" form:"page"`
	Limit     int    `json:"limit" form:"limit"`
}

type Recommendation struct {
	Ticker      string  `json:"ticker"`
	Company     string  `json:"company"`
	Score       float64 `json:"score"`
	Reason      string  `json:"reason"`
	TargetPrice string  `json:"target_price"`
	CurrentRating string `json:"current_rating"`
	Confidence  float64 `json:"confidence"`
}