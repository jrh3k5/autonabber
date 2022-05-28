package ynab

import "jrh3k5/autonabber/client/ynab/model"

type Client interface {
	GetBudgets() ([]*model.Budget, error)

	GetCategories(*model.Budget) ([]*model.BudgetCategoryGroup, error)
}
