package http

import (
	"encoding/json"
	"fmt"
	"jrh3k5/autonabber/client/ynab/model"
)

func (c *Client) GetBudgets() ([]*model.Budget, error) {
	httpResponse, err := c.Get("/budgets")
	if err != nil {
		return nil, fmt.Errorf("failed to invoke request for /budgets: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		errorContainer, err := handleError(httpResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to extract error details for request to /budgets: %w", err)
		}
		return nil, errorContainer
	}

	var response *budgetsResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response from /budgets: %w", err)
	}

	var budgets []*model.Budget
	if response.Data != nil {
		for _, budgetDetail := range response.Data.Budgets {
			budgets = append(budgets, &model.Budget{
				ID:   budgetDetail.ID,
				Name: budgetDetail.Name,
			})
		}
	}
	return budgets, nil
}
