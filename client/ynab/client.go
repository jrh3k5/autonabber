package ynab

import "jrh3k5/autonabber/client/ynab/model"

// Client is a definition of interactions that can be made with YNAB
type Client interface {
	// GetBudgets gets all currently-available budgets
	GetBudgets() ([]*model.Budget, error)

	// GetCategories gets all of the categories (in their groups) for the given budget
	GetCategories(budget *model.Budget) ([]*model.BudgetCategoryGroup, error)

	// SetBudget sets the budgeted amount for the given budget and category
	SetBudget(budget *model.Budget, category *model.BudgetCategory, newDollars int64, newCents int16) error
}
