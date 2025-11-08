package reward

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const idempotencyTTL = time.Hour

type IdempotencyService struct {
	cache *redis.Client
}

func NewIdempotencyService(client *redis.Client) *IdempotencyService {
	return &IdempotencyService{cache: client}
}

func (s *IdempotencyService) CheckOrSet(ctx context.Context, key string, response any) (bool, error) {
	val, err := s.cache.Get(ctx, key).Result()
	if err == nil && val != "" {
		return true, nil
	}

	if response != nil {
		jsonResp, _ := json.Marshal(response)
		_ = s.cache.Set(ctx, key, jsonResp, idempotencyTTL).Err()
	}

	return false, nil
}