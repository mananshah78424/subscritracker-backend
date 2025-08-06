package validator

import (
	"github.com/labstack/echo/v4"
)

type SubscriptionChannelRequest struct {
	ChannelName        string `json:"channel_name" form:"channel_name" validate:"required"`
	ChannelURL         string `json:"channel_url" form:"channel_url" validate:"required"`
	ChannelDescription string `json:"channel_description" form:"channel_description" validate:"required"`
	ChannelImageURL    string `json:"channel_image_url" form:"channel_image_url" validate:"required"`
}

func ValidateSubscriptionChannelRequest(c echo.Context) (*SubscriptionChannelRequest, error) {
	var request SubscriptionChannelRequest
	if err := c.Bind(&request); err != nil {
		return nil, err
	}

	if err := IsValidString(request.ChannelName, "channel_name"); err != nil {
		return nil, err
	}

	if err := IsValidString(request.ChannelURL, "channel_url"); err != nil {
		return nil, err
	}

	if err := IsValidString(request.ChannelDescription, "channel_description"); err != nil {
		return nil, err
	}

	if err := IsValidString(request.ChannelImageURL, "channel_image_url"); err != nil {
		return nil, err
	}

	return &request, nil
}
