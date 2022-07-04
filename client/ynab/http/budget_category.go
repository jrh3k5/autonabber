package http

import (
	"fmt"
	"time"

	"github.com/jrh3k5/autonabber/client/ynab/model"
)

func (c *Client) GetCategories(budget *model.Budget) ([]*model.BudgetCategoryGroup, error) {
	cacheKey := buildCategoryGroupsCacheKey(budget)

	if cached, cacheFound := c.ynabCache.Get(cacheKey); cacheFound {
		return cached.([]*model.BudgetCategoryGroup), nil
	}

	var budgetCategoriesResponse *budgetCategoriesResponse
	if err := c.GetJSON(fmt.Sprintf("/budgets/%s/categories", budget.ID), &budgetCategoriesResponse); err != nil {
		return nil, fmt.Errorf("failed to get categories for budget ID '%s': %w", budget.ID, err)
	}

	var budgetCategoryGroups []*model.BudgetCategoryGroup
	if budgetCategoriesResponse.Data != nil {
		for _, responseBudgetCategoryGroup := range budgetCategoriesResponse.Data.CategoryGroups {
			budgetCategoryGroup := &model.BudgetCategoryGroup{
				Name:   responseBudgetCategoryGroup.Name,
				Hidden: responseBudgetCategoryGroup.Hidden,
			}
			for _, responseBudgetCategory := range responseBudgetCategoryGroup.Categories {
				budgetedMillidollars := responseBudgetCategory.Budgeted % 1000
				budgetedDollars := (responseBudgetCategory.Budgeted - budgetedMillidollars) / 1000
				availableMillidollars := responseBudgetCategory.Balance % 1000
				availableDollars := (responseBudgetCategory.Balance - availableMillidollars) / 1000
				budgetCategory := &model.BudgetCategory{
					ID:               responseBudgetCategory.ID,
					Name:             responseBudgetCategory.Name,
					BudgetedDollars:  budgetedDollars,
					BudgetedCents:    int16(budgetedMillidollars / 10),
					AvailableDollars: availableDollars,
					AvailableCents:   int16(availableMillidollars / 10),
					Hidden:           responseBudgetCategory.Hidden,
				}
				budgetCategoryGroup.Categories = append(budgetCategoryGroup.Categories, budgetCategory)
			}
			budgetCategoryGroups = append(budgetCategoryGroups, budgetCategoryGroup)
		}
	}

	c.ynabCache.Add(cacheKey, budgetCategoryGroups, time.Hour)

	return budgetCategoryGroups, nil
}

func buildCategoryGroupsCacheKey(budget *model.Budget) string {
	return fmt.Sprintf("budgets/%s/categoryGroups", budget.ID)
}
