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
				budgetedDollars, budgetedCents := ParseMillidollars(responseBudgetCategory.Budgeted)
				availableDollars, availableCents := ParseMillidollars(responseBudgetCategory.Balance)
				budgetCategory := &model.BudgetCategory{
					ID:               responseBudgetCategory.ID,
					Name:             responseBudgetCategory.Name,
					BudgetedDollars:  budgetedDollars,
					BudgetedCents:    budgetedCents,
					AvailableDollars: availableDollars,
					AvailableCents:   availableCents,
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
