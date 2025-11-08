package price

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

// PriceService holds a redis client for caching stock prices
type PriceService struct {
	cache redisClient
}

// redisClient defines the methods we use from the redis.Client
type redisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// PriceResponse defines our JSON response model
type PriceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

// NewPriceService creates a new instance of PriceService
func NewPriceService(client *redis.Client) *PriceService {
	return &PriceService{
		cache: client,
	}
}

// GetStockPrice retrieves a stock price from cache or generates it randomly
func (p *PriceService) GetStockPrice(symbol string) (PriceResponse, error) {
	var resp PriceResponse
	ctx := context.Background()

	// 1. Try to fetch from cache
	cachedVal, err := p.cache.Get(ctx, symbol).Result()
	if err == nil && cachedVal != "" {
		if err := json.Unmarshal([]byte(cachedVal), &resp); err == nil {
			return resp, nil
		}
	}

	// 2. Generate a new random price between 1000 and 4000 INR
	price := 1000 + rand.Float64()*(4000-1000)
	resp = PriceResponse{
		Symbol: symbol,
		Price:  price,
	}

	// 3. Store in Redis for 10 minutes
	jsonData, _ := json.Marshal(resp)
	_ = p.cache.Set(ctx, symbol, jsonData, 10*time.Minute).Err()

	return resp, nil
}
