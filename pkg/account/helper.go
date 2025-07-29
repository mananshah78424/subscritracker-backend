package account

import (
	"context"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	"time"
)

func GetAccountById(app *application.App, id int) (*models.Account, error) {
	account := &models.Account{}

	err := app.Database.NewSelect().
		Model(account).
		Where("id = ?", id).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccountByGoogleID retrieves an account by Google ID
func GetAccountByGoogleID(app *application.App, googleID string) (*models.Account, error) {
	account := &models.Account{}

	err := app.Database.NewSelect().
		Model(account).
		Where("google_id = ?", googleID).
		Scan(context.Background())

	return account, err
}

func CreateAccount(app *application.App, account *models.Account) error {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	_, err := app.Database.NewInsert().
		Model(account).
		Exec(context.Background())

	return err
}

func UpdateAccountId(app *application.App, account *models.Account) error {
	account.UpdatedAt = time.Now()

	_, err := app.Database.NewUpdate().
		Model(account).
		Column("updated_at").
		Where("id = ?", account.ID).
		Exec(context.Background())

	return err
}

func DeleteAccount(app *application.App, account *models.Account) error {
	_, err := app.Database.NewDelete().
		Model(account).
		Where("id = ?", account.ID).
		Exec(context.Background())

	return err
}

func GetAccountStats(app *application.App, accountId int) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"total_subscriptions":  0,
		"active_subscriptions": 0,
		"monthly_spend":        0.00,
		"tier":                 "free",
		"features_used":        []string{"basic_tracking"},
	}

	return stats, nil
}
