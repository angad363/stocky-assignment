package users

import (
	"context"

	"github.com/angad363/stocky-assignment/internal/reward"
	"github.com/jmoiron/sqlx"
)

// UserService handles user creation and onboarding logic
type UserService struct {
	db        *sqlx.DB
	rewardSvc *reward.RewardService
}

func NewUserService(db *sqlx.DB, rewardSvc *reward.RewardService) *UserService {
	return &UserService{db: db, rewardSvc: rewardSvc}
}

// CreateUser inserts a new user and triggers an onboarding reward
func (s *UserService) CreateUser(ctx context.Context, name string) (User, reward.Reward, error) {
	var user User
	var rwd reward.Reward

	err := s.db.QueryRowContext(ctx,
		`INSERT INTO users (name) VALUES ($1)
		RETURNING id, name, created_at`,
		name,
	).Scan(&user.ID, &user.Name, &user.CreatedAt)
	if err != nil {
		return user, rwd, err
	}

	// Auto-reward new user with 1 share of a random stock
	rwdReq := reward.RewardRequest{
		UserID:   user.ID,
		Quantity: 1.0,
		Symbol:   "", // random stock
	}
	rwd, err = s.rewardSvc.CreateReward(ctx, rwdReq)
	if err != nil {
		return user, rwd, err
	}

	return user, rwd, nil
}