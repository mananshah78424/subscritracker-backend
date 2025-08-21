package account

import (
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ConvertAccountIdStringToInt(c echo.Context) (int, error) {
	userIDInterface := c.Get("user_id")
	var accountID int

	switch v := userIDInterface.(type) {
	case int:
		accountID = v
	case float64:
		accountID = int(v)
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			accountID = parsed
		} else {
			log.Printf("Failed to parse user_id from string: %v", err)
			return 0, err
		}
	default:
		return 0, fmt.Errorf("invalid user_id type: %T", userIDInterface)
	}

	return accountID, nil
}
