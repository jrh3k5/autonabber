package delta

import (
	"jrh3k5/autonabber/client/ynab/model"
	"jrh3k5/autonabber/input"
)

type BudgetCategoryDeltaGroup struct {
	Name           string
	CategoryDeltas []*BudgetCategoryDelta
}

type BudgetCategoryDelta struct {
	Name           string
	InitialDollars int64
	InitialCents   int16
	FinalDollars   int64
	FinalCents     int16
}

// CalculateDelta calculates the delta from the initial amount to the final amount.
// It returns the dollars and cents' difference.
// This currently only works with additive operations.
func (bcd *BudgetCategoryDelta) CalculateDelta() (int64, int16) {
	initialTotal := bcd.InitialDollars*100 + int64(bcd.InitialCents)
	finalTotal := bcd.FinalDollars*100 + int64(bcd.FinalCents)

	totalDiff := finalTotal - initialTotal
	totalDiffCents := totalDiff % 100
	totalDiffDollars := (totalDiff - totalDiffCents) / 100

	return totalDiffDollars, int16(totalDiffCents)
}

func NewDeltas(actual []*model.BudgetCategoryGroup, toApply *input.BudgetChange) ([]*BudgetCategoryDeltaGroup, error) {
	var deltaGroups []*BudgetCategoryDeltaGroup
	for _, actualGroup := range actual {
		deltaCategoryGroup := getCategoryGroupByName(actualGroup.Name, toApply)
		deltaGroup := &BudgetCategoryDeltaGroup{
			Name: actualGroup.Name,
		}
		for _, actualCategory := range actualGroup.Categories {
			deltaCategory := &BudgetCategoryDelta{
				Name:           actualCategory.Name,
				InitialDollars: actualCategory.BudgetedDollars,
				InitialCents:   actualCategory.BudgetedCents,
				FinalDollars:   actualCategory.BudgetedDollars,
				FinalCents:     actualCategory.BudgetedCents,
			}

			if deltaCategoryGroup != nil {
				if change := getChangeByName(actualCategory.Name, deltaCategoryGroup.Changes); change != nil {
					newDollars, newCents := change.ApplyDelta(actualCategory.BudgetedDollars, actualCategory.BudgetedCents)
					deltaCategory.FinalDollars = newDollars
					deltaCategory.FinalCents = newCents
				}
			}

			deltaGroup.CategoryDeltas = append(deltaGroup.CategoryDeltas, deltaCategory)
		}
		deltaGroups = append(deltaGroups, deltaGroup)
	}
	return deltaGroups, nil
}

func getCategoryGroupByName(name string, changes *input.BudgetChange) *input.BudgetCategoryGroup {
	for _, group := range changes.CategoryGroups {
		if name == group.Name {
			return group
		}
	}

	return nil
}

func getChangeByName(name string, changes []*input.BudgetCategoryChange) *input.BudgetCategoryChange {
	for _, change := range changes {
		if name == change.Name {
			return change
		}
	}

	return nil
}
