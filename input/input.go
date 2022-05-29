package input

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func ParseInputFile(filePath string) (*BudgetChanges, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}

	var serBudget *serializedBudgetChanges
	if err := yaml.Unmarshal(fileBytes, &serBudget); err != nil {
		return nil, fmt.Errorf("failed to parse YAML from file '%s': %w", filePath, err)
	}

	budgetChanges := &BudgetChanges{}
	for _, serChange := range serBudget.Changes {
		change := &BudgetChange{
			Name: serChange.Name,
		}
		for _, serGroup := range serChange.CategoryGroups {
			categoryGroup := &BudgetCategoryGroup{
				Name: serGroup.Name,
			}
			for _, serCategory := range serGroup.Categories {
				category, err := NewBudgetCategoryChange(serCategory.Name, serCategory.Change)
				if err != nil {
					return nil, fmt.Errorf("unable to create change for budget category '%s': %w", serCategory.Name, err)
				}
				categoryGroup.Changes = append(categoryGroup.Changes, category)
			}
			change.CategoryGroups = append(change.CategoryGroups, categoryGroup)
		}
		budgetChanges.Changes = append(budgetChanges.Changes, change)
	}
	return budgetChanges, nil
}
