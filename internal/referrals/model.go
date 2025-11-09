package referral

import "time"

type Referral struct {
	ID         int       `db:"id" json:"id"`
	ReferrerID int       `db:"referrer_id" json:"referrer_id"`
	FriendName string    `db:"friend_name" json:"friend_name"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
