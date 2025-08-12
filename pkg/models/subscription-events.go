package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Subscription_Event struct {
	bun.BaseModel         `bun:"subscription_events"`
	ID                    int       `json:"id" bun:",autoincrement"`
	SubscriptionDetailsID int       `json:"subscription_details_id"`
	AccountID             int       `json:"account_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
