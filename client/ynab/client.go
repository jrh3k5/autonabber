package ynab

import (
	"github.com/jrh3k5/autonabber/client/ynab/model"
)

// Client is a definition of interactions that can be made with YNAB
type Client interface {
	// GetAverageSpent gets the average amount spent each month in the given category over the past monthLookback months
	// It returns the dollars and cents of the average
	GetMonthlyAverageSpent(budget *model.Budget, category *model.BudgetCategory, monthLookback int) (int64, int16, error)

	// GetBudgets gets all currently-available budgets
	GetBudgets() ([]*model.Budget, error)

	// GetCategories gets all of the categories (in their groups) for the given budget
	GetCategories(budget *model.Budget) ([]*model.BudgetCategoryGroup, error)

	// SetBudget sets the budgeted amount for the given budget and category
	SetBudget(budget *model.Budget, category *model.BudgetCategory, newDollars int64, newCents int16) error
}
