package reward

import "time"

type Reward struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	StockSymbol string    `db:"stock_symbol" json:"stock_symbol"`
	Quantity    float64   `db:"quantity" json:"quantity"`
	RewardedAt  time.Time `db:"rewarded_at" json:"rewarded_at"`
}

type RewardRequest struct {
	UserID   int     `json:"user_id"`
	Symbol   string  `json:"symbol,omitempty"`
	Quantity float64 `json:"quantity"`
}
