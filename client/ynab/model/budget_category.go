package model

import (
	"fmt"

	"github.com/jrh3k5/autonabber/format"
)

type BudgetCategoryGroup struct {
	Name       string
	Categories []*BudgetCategory
	Hidden     bool
}

// AllCategoriesHidden determines if all of the categories within this group are hidden
func (bcg *BudgetCategoryGroup) AllCategoriesHidden() bool {
	for _, category := range bcg.Categories {
		if !category.Hidden {
			return false
		}
	}
	return true
}

type BudgetCategory struct {
	ID               string
	Name             string
	BudgetedDollars  int64
	BudgetedCents    int16
	AvailableDollars int64
	AvailableCents   int16
	Hidden           bool
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

func PrintBudgetCategoryGroups(budgetCategoryGroups []*BudgetCategoryGroup, printHiddenCategories bool) {
	var indent int
	for _, budgetCategoryGroup := range budgetCategoryGroups {
		// Don't print out the internal master category - doesn't seem terribly helpful
		if "Internal Master Category" == budgetCategoryGroup.Name {
			continue
		}

		if (budgetCategoryGroup.Hidden || budgetCategoryGroup.AllCategoriesHidden()) && !printHiddenCategories {
			continue
		}

		fmt.Println(budgetCategoryGroup.Name)
		indent += 2
		for _, budgetCategory := range budgetCategoryGroup.Categories {
			if budgetCategory.Hidden && !printHiddenCategories {
				continue
			}

			formattedAvailable := format.FormatUSD(budgetCategory.AvailableDollars, budgetCategory.AvailableCents)
			formattedBudgeted := format.FormatUSD(budgetCategory.BudgetedDollars, budgetCategory.BudgetedCents)
			fmt.Printf("  %s (budgeted: %s; total available: %s)\n", budgetCategory.Name, formattedBudgeted, formattedAvailable)
		}
	}
}
