package delta

import (
	"context"
	"fmt"

	"github.com/jrh3k5/autonabber/client/ynab"
	"github.com/jrh3k5/autonabber/client/ynab/model"
	"github.com/jrh3k5/autonabber/format"
	"go.uber.org/zap"
)

// Applicator defines a means to apply changes to a budget
type Applicator interface {
	// ApplyChanges applies the given changes, return the total dollars and cents that were applied
	ApplyChanges(ctx context.Context, budget *model.Budget, deltas []*BudgetCategoryDeltaGroup) (int64, int16, error)
}

// ClientApplicator is an implementation of Applicator backed by an instance of ynab.Client
type ClientApplicator struct {
	logger *zap.Logger
	client ynab.Client
}

func NewClientApplicator(logger *zap.Logger, client ynab.Client) *ClientApplicator {
	return &ClientApplicator{
		logger: logger,
		client: client,
	}
}

func (c *ClientApplicator) ApplyChanges(ctx context.Context, budget *model.Budget, deltas []*BudgetCategoryDeltaGroup) (int64, int16, error) {
	var nonZeroChanges []*BudgetCategoryDelta

	for _, delta := range deltas {
		for _, change := range delta.CategoryDeltas {
			if change.HasChanges() {
				nonZeroChanges = append(nonZeroChanges, change)
			}
		}
	}

	for changeIndex, change := range nonZeroChanges {
		c.logger.Sugar().Infof("Applying change %d of %d", changeIndex+1, len(nonZeroChanges))
		if err := c.client.SetBudget(budget, change.BudgetCategory, change.FinalBudgetDollars, change.FinalBudgetCents); err != nil {
			formattedFinal := format.FormatUSD(change.FinalDollars, change.FinalCents)
			return 0, 0, fmt.Errorf("failed to set budget category '%s' under budget '%s' to %s: %w", change.BudgetCategory.Name, budget.Name, formattedFinal, err)
		}
	}

	deltaDollars, deltaCents := SumChanges(nonZeroChanges)
	return deltaDollars, deltaCents, nil
}
