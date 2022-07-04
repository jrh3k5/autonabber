package http

import (
	"fmt"

	"github.com/jrh3k5/autonabber/client/ynab/model"
)

func (c *Client) GetReadyToAssign(budget *model.Budget) (int64, int16, error) {
	requestPath := fmt.Sprintf("/budgets/%s/months/current", budget.ID)
	var monthData *budgetMonthData
	err := c.GetJSON(requestPath, &monthData)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to invoke request for %s: %w", requestPath, err)
	}

	incomeMillis := monthData.Data.Month.ToBeBudgeted
	incomeDollars, incomeCents := ParseMillidollars(incomeMillis)
	return incomeDollars, incomeCents, nil
}
