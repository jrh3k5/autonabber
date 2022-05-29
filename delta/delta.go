package delta

import (
	"fmt"
	"jrh3k5/autonabber/client/ynab/model"
	"jrh3k5/autonabber/input"
)

type BudgetCategoryDeltaGroup struct {
	Name           string
	CategoryDeltas []*BudgetCategoryDelta
}

type BudgetCategoryDelta struct {
	BudgetCategory *model.BudgetCategory
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

// HasChanges determines if this delta has any actual changes to be applied
func (bcd *BudgetCategoryDelta) HasChanges() bool {
	dollars, cents := bcd.CalculateDelta()
	return dollars != 0 || cents != 0
}

// SumChanges calculates the sum of all changes across the given groups
// It returns the dollars and cents of the changes
func SumChanges(deltas []*BudgetCategoryDelta) (int64, int16) {
	diffDollars := int64(0)
	diffCents := int64(0)

	for _, category := range deltas {
		catDollars, catCents := category.CalculateDelta()
		diffDollars += int64(catDollars)
		diffCents += int64(catCents)
	}

	totalDiff := diffDollars*100 + diffCents
	finalCents := totalDiff % 100
	finalDollars := (totalDiff - finalCents) / 100

	return finalDollars, int16(finalCents)
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
				BudgetCategory: actualCategory,
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

func PrintDeltas(groups []*BudgetCategoryDeltaGroup) {
	for _, group := range groups {
		var nonZeroChanges []*BudgetCategoryDelta
		for _, change := range group.CategoryDeltas {
			if change.HasChanges() {
				nonZeroChanges = append(nonZeroChanges, change)
			}
		}

		if len(nonZeroChanges) == 0 {
			// Don't print any groups without changes
			continue
		}

		fmt.Printf("%s\n", group.Name)
		for _, change := range nonZeroChanges {
			deltaDollars, deltaCents := change.CalculateDelta()
			fmt.Printf("  %s: $%d.%02d => $%d.%02d (+$%d.%02d)\n", change.BudgetCategory.Name, change.InitialDollars, change.InitialCents, change.FinalDollars, change.FinalCents, deltaDollars, deltaCents)
		}
	}
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
