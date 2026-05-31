package ynab

import (
	"fmt"
	"sync"

	"codeberg.org/jrh3k5/autonabber/client/ynab/model"
)

type monthlyAverageSpent struct {
	dollars int64
	cents   int16
}

type readyToAssignValue struct {
	dollars int64
	cents   int16
}

type InMemoryClient struct {
	mu sync.RWMutex

	budgets       map[string]*model.Budget
	categories    map[string][]*model.BudgetCategoryGroup
	averages      map[string]map[string]map[int]*monthlyAverageSpent
	readyToAssign map[string]*readyToAssignValue
}

func NewInMemoryClient() *InMemoryClient {
	return &InMemoryClient{
		budgets:       make(map[string]*model.Budget),
		categories:    make(map[string][]*model.BudgetCategoryGroup),
		averages:      make(map[string]map[string]map[int]*monthlyAverageSpent),
		readyToAssign: make(map[string]*readyToAssignValue),
	}
}

func (c *InMemoryClient) GetMonthlyAverageSpent(budget *model.Budget, category *model.BudgetCategory, monthLookback int) (int64, int16, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	budgetAverages, ok := c.averages[budget.ID]
	if !ok {
		return 0, 0, fmt.Errorf("no averages found for budget '%s'", budget.ID)
	}

	categoryAverages, ok := budgetAverages[category.ID]
	if !ok {
		return 0, 0, fmt.Errorf("no averages found for category '%s'", category.ID)
	}

	average, ok := categoryAverages[monthLookback]
	if !ok {
		return 0, 0, fmt.Errorf("no average found for month lookback %d", monthLookback)
	}

	return average.dollars, average.cents, nil
}

func (c *InMemoryClient) GetBudgets() ([]*model.Budget, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	budgets := make([]*model.Budget, 0, len(c.budgets))
	for _, budget := range c.budgets {
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (c *InMemoryClient) GetCategories(budget *model.Budget) ([]*model.BudgetCategoryGroup, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	categories, ok := c.categories[budget.ID]
	if !ok {
		return nil, fmt.Errorf("no categories found for budget '%s'", budget.ID)
	}

	return categories, nil
}

func (c *InMemoryClient) GetReadyToAssign(budget *model.Budget) (int64, int16, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.readyToAssign[budget.ID]
	if !ok {
		return 0, 0, fmt.Errorf("no ready-to-assign value found for budget '%s'", budget.ID)
	}

	return value.dollars, value.cents, nil
}

func (c *InMemoryClient) SetBudget(budget *model.Budget, category *model.BudgetCategory, newDollars int64, newCents int16) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	groups, ok := c.categories[budget.ID]
	if !ok {
		return fmt.Errorf("no categories found for budget '%s'", budget.ID)
	}

	for _, group := range groups {
		for _, cat := range group.Categories {
			if cat.ID == category.ID {
				cat.BudgetedDollars = newDollars
				cat.BudgetedCents = newCents
				return nil
			}
		}
	}

	return fmt.Errorf("category '%s' not found in budget '%s'", category.ID, budget.ID)
}

func (c *InMemoryClient) SetBudgets(budgets []*model.Budget) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.budgets = make(map[string]*model.Budget, len(budgets))
	for _, budget := range budgets {
		c.budgets[budget.ID] = budget
	}
}

func (c *InMemoryClient) SetCategories(budgetID string, groups []*model.BudgetCategoryGroup) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.categories[budgetID] = groups
}

func (c *InMemoryClient) SetMonthlyAverageSpent(budgetID string, categoryID string, monthLookback int, dollars int64, cents int16) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.averages[budgetID] == nil {
		c.averages[budgetID] = make(map[string]map[int]*monthlyAverageSpent)
	}

	if c.averages[budgetID][categoryID] == nil {
		c.averages[budgetID][categoryID] = make(map[int]*monthlyAverageSpent)
	}

	c.averages[budgetID][categoryID][monthLookback] = &monthlyAverageSpent{
		dollars: dollars,
		cents:   cents,
	}
}

func (c *InMemoryClient) SetReadyToAssign(budgetID string, dollars int64, cents int16) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.readyToAssign[budgetID] = &readyToAssignValue{
		dollars: dollars,
		cents:   cents,
	}
}
