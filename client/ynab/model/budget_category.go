package model

import (
	"fmt"
)

type BudgetCategoryGroup struct {
	Name       string
	Categories []*BudgetCategory
}

type BudgetCategory struct {
	ID              string
	Name            string
	BudgetedDollars int64
	BudgetedCents   int16
}

func PrintBudgetCategoryGroups(budgetCategoryGroups []*BudgetCategoryGroup) {
	var indent int
	for _, budgetCategoryGroup := range budgetCategoryGroups {
		indent = 0
		fmt.Println(buildIndentation(indent) + budgetCategoryGroup.Name)
		indent += 2
		for _, budgetCategory := range budgetCategoryGroup.Categories {
			fmt.Printf("%s%s ($%d.%02d)\n", buildIndentation(indent), budgetCategory.Name, budgetCategory.BudgetedDollars, budgetCategory.BudgetedCents)
		}
	}
}
