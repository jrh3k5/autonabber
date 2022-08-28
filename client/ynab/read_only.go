package ynab

import (
	"github.com/jrh3k5/autonabber/client/ynab/model"
	"go.uber.org/zap"
)

type readOnlyClient struct {
	logger     *zap.SugaredLogger
	liveClient Client
}

func NewReadOnlyClient(logger *zap.SugaredLogger, liveClient Client) Client {
	return &readOnlyClient{
		logger:     logger,
		liveClient: liveClient,
	}
}

func (c *readOnlyClient) GetMonthlyAverageSpent(budget *model.Budget, category *model.BudgetCategory, monthLookback int) (int64, int16, error) {
	return c.liveClient.GetMonthlyAverageSpent(budget, category, monthLookback)
}

func (c *readOnlyClient) GetBudgets() ([]*model.Budget, error) {
	return c.liveClient.GetBudgets()
}

func (c *readOnlyClient) GetCategories(budget *model.Budget) ([]*model.BudgetCategoryGroup, error) {
	return c.liveClient.GetCategories(budget)
}

func (c *readOnlyClient) GetReadyToAssign(budget *model.Budget) (int64, int16, error) {
	return c.liveClient.GetReadyToAssign(budget)
}

func (c *readOnlyClient) SetBudget(budget *model.Budget, category *model.BudgetCategory, newDollars int64, newCents int16) error {
	c.logger.Infof("Client is read-only; would have, otherwise, applied %d new dollars and %d new cents to the budget '%s' (ID: '%s')", newDollars, newCents, budget.Name, budget.ID)
	return nil
}
