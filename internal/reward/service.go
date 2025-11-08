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