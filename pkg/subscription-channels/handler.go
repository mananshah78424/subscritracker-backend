package subscription_channels

import (
	"log"
	"net/http"
	"subscritracker/pkg/models"
	"subscritracker/pkg/validator"

	"github.com/labstack/echo/v4"
)

func GetSubscriptionChannelsHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		log.Println("No id found while searching for subscription")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	channel, err := GetChannelById(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channel)
}

func GetAllSubscriptionChannelsHandler(c echo.Context) error {
	channels, err := GetAllChannels(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channels)
}

func PostSubscriptionChannelsHandler(c echo.Context) error {
	// Validate the request and get the validated data
	request, err := validator.ValidateSubscriptionChannelRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Create the model from the validated request
	channel := models.Subscription_Channels{
		ChannelName:        request.ChannelName,
		ChannelURL:         request.ChannelURL,
		ChannelDescription: request.ChannelDescription,
		ChannelImageURL:    request.ChannelImageURL,
	}

	createdChannel, err := CreateChannel(c, channel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdChannel)
}
