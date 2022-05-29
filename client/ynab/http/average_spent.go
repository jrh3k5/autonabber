package http

import (
	"fmt"
	"jrh3k5/autonabber/client/ynab/model"
	"math"
	"time"
)

func (c *Client) GetMonthlyAverageSpent(budget *model.Budget, category *model.BudgetCategory, monthLookback int) (int64, int16, error) {
	cacheKey := fmt.Sprintf("/budgets/%s/categories/%s/transactions-avg-spent-%dm", budget.ID, category.ID, monthLookback)
	if cached, cacheFound := c.ynabCache.Get(cacheKey); cacheFound {
		totalAverage := cached.(int64)
		averageCents := totalAverage % 100
		averageDollars := (totalAverage - averageCents) / 100

		return averageDollars, int16(averageCents), nil
	}

	since := time.Now().AddDate(0, -1*monthLookback, 0)
	requestPath := fmt.Sprintf("/budgets/%s/categories/%s/transactions?since_date=%s", budget.ID, category.ID, since.Format("2006-01-02"))
	var response *transactionsContainer
	if err := c.GetJSON(requestPath, &response); err != nil {
		return 0, 0, fmt.Errorf("failed to invoke GET request to '%s': %w", requestPath, err)
	}

	var totalAmounts int64
	if response != nil && response.Data != nil {
		for _, transaction := range response.Data.Transactions {
			// Invert so that expenses (which are negative amounts) are added to the total of costs
			// while credits (which are positive amounts) are deducted from the overall amount spent
			totalAmounts += -1 * transaction.Amount
		}
	}

	// Because it's expressed in millidollars, split the amount by 10 to get to dollars and cents
	totalAmounts /= 10
	averageAmount := int64(math.Ceil(float64(totalAmounts) / float64(monthLookback)))
	averageCents := averageAmount % 100
	averageDollars := (averageAmount - averageCents) / 100

	c.ynabCache.Add(cacheKey, averageAmount, time.Hour)

	return averageDollars, int16(averageCents), nil
}
