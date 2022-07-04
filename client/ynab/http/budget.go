package http

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jrh3k5/autonabber/client/ynab/model"
	"github.com/jrh3k5/autonabber/format"
)

func (c *Client) GetBudgets() ([]*model.Budget, error) {
	cacheKey := "/budgets"

	if cached, cacheFound := c.ynabCache.Get(cacheKey); cacheFound {
		return cached.([]*model.Budget), nil
	}

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

	c.ynabCache.Add(cacheKey, budgets, time.Hour)

	return budgets, nil
}

func (c *Client) SetBudget(budget *model.Budget, category *model.BudgetCategory, newDollars int64, newCents int16) error {
	budgetedMillidollars := newDollars*1000 + int64(newCents)*100
	requestPath := fmt.Sprintf("/budgets/%s/months/current/categories/%s", budget.ID, category.ID)
	requestBody := &categoryPatchRequest{
		Category: &patchedCategory{
			Budgeted: budgetedMillidollars,
		},
	}

	if err := c.PatchJSON(requestPath, requestBody); err != nil {
		formattedNew := format.FormatUSD(newDollars, newCents)
		return fmt.Errorf("failed to update category '%s' in budget '%s' to %s: %w", category.Name, budget.Name, formattedNew, err)
	}

	// Because the category groups have been changed, evict the stale data from the cache
	categoryGroupsCacheKey := buildCategoryGroupsCacheKey(budget)
	c.ynabCache.Delete(categoryGroupsCacheKey)

	return nil
}
