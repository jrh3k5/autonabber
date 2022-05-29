package input

import (
	"fmt"

	"github.com/jrh3k5/autonabber/client/ynab"
	"github.com/jrh3k5/autonabber/client/ynab/model"
)

type BudgetChanges struct {
	Changes []*BudgetChange
}

type BudgetChange struct {
	Name           string
	CategoryGroups []*BudgetCategoryGroup
}

type BudgetCategoryGroup struct {
	Name    string
	Changes []*BudgetCategoryChange
}

func NewBudgetCategoryChange(name string, changeOperation string) (*BudgetCategoryChange, error) {
	var changeOpFunction changeOperationFn
	if addOpRegexp.MatchString(changeOperation) {
		changeOpFunction = addRawValue
	} else if averageMonthlySpentRegex.MatchString(changeOperation) {
		changeOpFunction = addMonthlyAverage
	} else {
		return nil, fmt.Errorf("unrecognized change operation for budget category '%s': %s", name, changeOperation)
	}

	return &BudgetCategoryChange{
		Name:                    name,
		changeOperation:         changeOperation,
		changeOperationFunction: changeOpFunction,
	}, nil
}

type BudgetCategoryChange struct {
	Name                    string
	changeOperation         string
	changeOperationFunction changeOperationFn
}

// ApplyDelta will apply the change described in the budget category change
// It returns the given dollars and cents after this change has been applied
func (bcc *BudgetCategoryChange) ApplyDelta(client ynab.Client, budget *model.Budget, budgetCategory *model.BudgetCategory, dollars int64, cents int16) (int64, int16, error) {
	if bcc.changeOperationFunction == nil {
		// Not graceful, but prior validation should prevent us from reaching this point
		panic(fmt.Sprintf("Unexpected operation in change for category '%s': %s", bcc.Name, bcc.changeOperation))
	}

	return bcc.changeOperationFunction(bcc, client, budget, budgetCategory, dollars, cents)
}
