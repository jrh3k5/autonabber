package http

import (
	"fmt"
	"jrh3k5/autonabber/client/ynab/model"
)

func (c *Client) GetCategories(budget *model.Budget) ([]*model.BudgetCategoryGroup, error) {
	var budgetCategoriesResponse *budgetCategoriesResponse
	if err := c.GetJSON(fmt.Sprintf("/budgets/%s/categories", budget.ID), &budgetCategoriesResponse); err != nil {
		return nil, fmt.Errorf("failed to get categories for budget ID '%s': %w", budget.ID, err)
	}

	var budgetCategoryGroups []*model.BudgetCategoryGroup
	if budgetCategoriesResponse.Data != nil {
		for _, responseBudgetCategoryGroup := range budgetCategoriesResponse.Data.CategoryGroups {
			budgetCategoryGroup := &model.BudgetCategoryGroup{
				Name: responseBudgetCategoryGroup.Name,
			}
			for _, responseBudgetCategory := range responseBudgetCategoryGroup.Categories {
				budgetedMillidollars := responseBudgetCategory.Budgeted % 1000
				budgetedDollars := (responseBudgetCategory.Budgeted - budgetedMillidollars) / 1000
				budgetCategory := &model.BudgetCategory{
					Name:            responseBudgetCategory.Name,
					BudgetedDollars: budgetedDollars,
					BudgetedCents:   int16(budgetedMillidollars / 10),
				}
				budgetCategoryGroup.Categories = append(budgetCategoryGroup.Categories, budgetCategory)
			}
			budgetCategoryGroups = append(budgetCategoryGroups, budgetCategoryGroup)
		}
	}
	return budgetCategoryGroups, nil
}
