package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Subscription_Channels struct {
	bun.BaseModel      `bun:"subscription_channels"`
	ID                 int       `bun:"id,pk,autoincrement" json:"id"`
	ChannelName        string    `bun:"channel_name,unique" json:"channel_name"`
	ChannelURL         string    `bun:"channel_url,unique" json:"channel_url"`
	ChannelDescription string    `bun:"channel_description" json:"channel_description"`
	ChannelImageURL    string    `bun:"channel_image_url" json:"channel_image_url"`
	ChannelType        string    `bun:"channel_type" json:"channel_type"`
	ChannelStatus      string    `bun:"channel_status" json:"channel_status"`
	CreatedAt          time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt          time.Time `bun:"updated_at" json:"updated_at"`
}
