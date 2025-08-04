package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel     `bun:"account"`
	ID                int                    `bun:"id,pk,autoincrement" json:"id"`
	GoogleID          *string                `bun:"google_id,unique" json:"google_id"`
	Email             string                 `bun:"email,unique" json:"email"`
	Name              string                 `bun:"name" json:"name"`
	GivenName         string                 `bun:"given_name" json:"given_name"`
	FamilyName        string                 `bun:"family_name" json:"family_name"`
	PictureURL        string                 `bun:"picture_url" json:"picture_url"`
	EmailVerified     bool                   `bun:"email_verified" json:"email_verified"`
	PasswordHash      string                 `bun:"password_hash" json:"password_hash"`
	VerificationToken string                 `bun:"verification_token" json:"verification_token"`
	ResetToken        string                 `bun:"reset_token" json:"reset_token"`
	ResetTokenExpires time.Time              `bun:"reset_token_expires" json:"reset_token_expires"`
	Tier              string                 `bun:"tier" json:"tier"`
	Status            string                 `bun:"status" json:"status"`
	Features          map[string]interface{} `bun:"features" json:"features"`
	SubscriptionCount int                    `bun:"subscription_count" json:"subscription_count"`
	LastLoginAt       time.Time              `bun:"last_login_at" json:"last_login_at"`
	CreatedAt         time.Time              `bun:"created_at" json:"created_at"`
	UpdatedAt         time.Time              `bun:"updated_at" json:"updated_at"`
}
