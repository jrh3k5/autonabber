package undo

import (
	"context"

	"github.com/jrh3k5/autonabber/client/ynab/model"
	"github.com/jrh3k5/autonabber/delta"
	"github.com/jrh3k5/autonabber/input"
)

// Generator defines a means to generate an undo
type Generator interface {
	// GenerateUndo builds an undo data object to undo the given deltas derived from the given change applied against the given budget
	GenerateUndo(ctx context.Context, budget *model.Budget, change *input.BudgetChange, deltas []*delta.BudgetCategoryDeltaGroup) (*BudgetChangesUndo, error)
}

type GeneratorImpl struct {
}

func NewGeneratorImpl() *GeneratorImpl {
	return &GeneratorImpl{}
}

func (*GeneratorImpl) GenerateUndo(ctx context.Context, budget *model.Budget, change *input.BudgetChange, deltas []*delta.BudgetCategoryDeltaGroup) (*BudgetChangesUndo, error) {

}
