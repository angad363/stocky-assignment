package reward

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/angad363/stocky-assignment/internal/price"
	"github.com/jmoiron/sqlx"
)

type RewardService struct {
	db       *sqlx.DB
	priceSvc *price.PriceService
}

func NewRewardService(db *sqlx.DB, priceSvc *price.PriceService) *RewardService {
	return &RewardService{db: db, priceSvc: priceSvc}
}

func (s *RewardService) CreateReward(ctx context.Context, req RewardRequest) (Reward, error) {
	var reward Reward

	symbol := req.Symbol
	if symbol == "" {
		stocks := []string{"RELIANCE", "TCS", "INFY", "HDFC", "ICICIBANK"}
		symbol = stocks[rand.Intn(len(stocks))]
	}

	// we call the price service just to simulate reward logic,
	// though we don’t store the price in DB
	_, err := s.priceSvc.GetStockPrice(symbol)
	if err != nil {
		return reward, err
	}

	reward = Reward{
		UserID:      req.UserID,
		StockSymbol: symbol,
		Quantity:    req.Quantity,
		RewardedAt:  time.Now(),
	}

	query := `
		INSERT INTO rewards (user_id, stock_symbol, quantity, rewarded_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = s.db.QueryRowContext(ctx, query,
		reward.UserID,
		reward.StockSymbol,
		reward.Quantity,
		reward.RewardedAt,
	).Scan(&reward.ID)


	if err != nil {
		fmt.Println("❌ Insert error:", err)
		return reward, err
	}

	return reward, err
}

func (s *RewardService) GetTodayRewards(ctx context.Context, userID int) ([]Reward, error) {
	rewards := []Reward{}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT id, user_id, stock_symbol, quantity, rewarded_at
		FROM rewards
		WHERE user_id = $1
		  AND rewarded_at >= $2
		  AND rewarded_at < $3
		ORDER BY rewarded_at DESC
	`

	err := s.db.SelectContext(ctx, &rewards, query, userID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

func (s *RewardService) GetHistoricalINR(ctx context.Context, userID int) ([]HistoricalINR, error) {
	rows, err := s.db.QueryxContext(ctx, `
		SELECT stock_symbol, quantity, rewarded_at
		FROM rewards
		WHERE user_id = $1
		  AND rewarded_at < CURRENT_DATE
		ORDER BY rewarded_at ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := []Reward{}
	for rows.Next() {
		var r Reward
		if err := rows.StructScan(&r); err != nil {
			return nil, err
		}
		rewards = append(rewards, r)
	}

	// Group total INR per date
	dateTotals := make(map[string]float64)

	for _, r := range rewards {
		price, err := s.priceSvc.GetStockPrice(r.StockSymbol)
		if err != nil {
			continue
		}
		inrValue := r.Quantity * price.Price

		dateKey := r.RewardedAt.Format("2006-01-02")
		dateTotals[dateKey] += inrValue
	}

	var historical []HistoricalINR
	for date, total := range dateTotals {
		historical = append(historical, HistoricalINR{
			Date:     date,
			TotalINR: total,
		})
	}

	return historical, nil
}

func (s *RewardService) GetUserStats(ctx context.Context, userID int) (map[string]float64, float64, error) {
	todayQuery := `
		SELECT stock_symbol, SUM(quantity) AS total_quantity
		FROM rewards
		WHERE user_id = $1
		AND DATE(rewarded_at AT TIME ZONE 'Asia/Kolkata') = CURRENT_DATE
		GROUP BY stock_symbol
	`
	todayRows, err := s.db.QueryxContext(ctx, todayQuery, userID)
	if err != nil {
		return nil, 0, err
	}
	defer todayRows.Close()

	todaySummary := make(map[string]float64)
	for todayRows.Next() {
		var symbol string
		var qty float64
		if err := todayRows.Scan(&symbol, &qty); err == nil {
			todaySummary[symbol] = qty
		}
	}

	holdingsQuery := `
		SELECT stock_symbol, SUM(quantity) AS total_quantity
		FROM rewards
		WHERE user_id = $1
		GROUP BY stock_symbol
	`
	holdRows, err := s.db.QueryxContext(ctx, holdingsQuery, userID)
	if err != nil {
		return todaySummary, 0, err
	}
	defer holdRows.Close()

	totalValue := 0.0
	for holdRows.Next() {
		var symbol string
		var qty float64
		if err := holdRows.Scan(&symbol, &qty); err == nil {
			priceResp, err := s.priceSvc.GetStockPrice(symbol)
			if err != nil {
				continue
			}
			// handle rounding for INR precision
			inr := qty * priceResp.Price
			inr = float64(int(inr*100+0.5)) / 100 // round to 2 decimals
			totalValue += inr
		}
	}

	return todaySummary, totalValue, nil
}

