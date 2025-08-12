package subscription_channels

import (
	"context"
	"log"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	"time"

	"github.com/labstack/echo/v4"
)

func GetChannelById(c echo.Context, id string) (*models.Subscription_Channels, error) {
	channel := &models.Subscription_Channels{}
	app := c.Get("app").(*application.App)
	err := app.Database.NewSelect().
		Model(channel).
		Where("id = ?", id).
		Scan(context.Background())

	if err != nil {
		log.Println("Error getting channel by id: ", err)
		return nil, err
	}

	return channel, nil

}

func GetAllChannels(c echo.Context) ([]*models.Subscription_Channels, error) {
	channels := []*models.Subscription_Channels{}
	app := c.Get("app").(*application.App)
	err := app.Database.NewSelect().
		Model(&channels).
		Scan(context.Background())

	if err != nil {
		log.Println("Error getting all channels: ", err)
		return nil, err
	}

	return channels, nil
}

func CreateChannel(c echo.Context, channel models.Subscription_Channels) (*models.Subscription_Channels, error) {
	app := c.Get("app").(*application.App)

	// Set timestamps
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	// Set default values if not provided
	if channel.ChannelStatus == "" {
		channel.ChannelStatus = "active"
	}
	if channel.ChannelType == "" {
		channel.ChannelType = "streaming"
	}

	_, err := app.Database.NewInsert().
		Model(&channel).
		Exec(context.Background())

	if err != nil {
		log.Println("Error creating channel: ", err)
		return nil, err
	}

	return &channel, nil
}
