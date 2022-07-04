package model

import (
	"fmt"

	"github.com/jrh3k5/autonabber/format"
)

type BudgetCategoryGroup struct {
	Name       string
	Categories []*BudgetCategory
}

type BudgetCategory struct {
	ID               string
	Name             string
	BudgetedDollars  int64
	BudgetedCents    int16
	AvailableDollars int64
	AvailableCents   int16
}

// GetReadyToAssign reads the "Ready to Assign" balance, if found, from the given groups
func GetReadyToAssign(groups []*BudgetCategoryGroup) (int64, int16) {
	for _, group := range groups {
		if group.Name == "Internal Master Category" {
			for _, category := range group.Categories {
				if category.Name == "Inflow: Ready to Assign" {
					return category.BudgetedDollars, category.BudgetedCents
				}
			}
		}
	}

	return 0, 0
}

func PrintBudgetCategoryGroups(budgetCategoryGroups []*BudgetCategoryGroup) {
	var indent int
	for _, budgetCategoryGroup := range budgetCategoryGroups {
		// Don't print out the internal master category - doesn't seem terribly helpful
		if "Internal Master Category" == budgetCategoryGroup.Name {
			continue
		}
		fmt.Println(budgetCategoryGroup.Name)
		indent += 2
		for _, budgetCategory := range budgetCategoryGroup.Categories {
			formattedAvailable := format.FormatUSD(budgetCategory.AvailableDollars, budgetCategory.AvailableCents)
			formattedBudgeted := format.FormatUSD(budgetCategory.BudgetedDollars, budgetCategory.BudgetedCents)
			fmt.Printf("  %s (budgeted: %s; total available: %s)\n", budgetCategory.Name, formattedBudgeted, formattedAvailable)
		}
	}
}
