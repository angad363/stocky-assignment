package referral

import (
	"context"
	"time"

	"github.com/angad363/stocky-assignment/internal/reward"
	"github.com/jmoiron/sqlx"
)

type ReferralService struct {
	db        *sqlx.DB
	rewardSvc *reward.RewardService
}

func NewReferralService(db *sqlx.DB, rewardSvc *reward.RewardService) *ReferralService {
	return &ReferralService{db: db, rewardSvc: rewardSvc}
}

func (s *ReferralService) CreateReferral(ctx context.Context, referrerID int, friendName string) (Referral, reward.Reward, error) {
	var ref Referral
	var rwd reward.Reward

	err := s.db.QueryRowContext(ctx, `
		INSERT INTO referrals (referrer_id, friend_name, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, referrer_id, friend_name, created_at
	`, referrerID, friendName, time.Now()).Scan(
		&ref.ID, &ref.ReferrerID, &ref.FriendName, &ref.CreatedAt,
	)
	if err != nil {
		return ref, rwd, err
	}

	// Give a random stock reward for successful referral
	rwdReq := reward.RewardRequest{
		UserID:   ref.ReferrerID,
		Quantity: 1.0,
		Symbol:   "",
	}
	rwd, err = s.rewardSvc.CreateReward(ctx, rwdReq)
	if err != nil {
		return ref, rwd, err
	}

	return ref, rwd, nil
}