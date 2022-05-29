package main

import (
	"errors"
	"fmt"
	"jrh3k5/autonabber/args"
	"jrh3k5/autonabber/client/ynab"
	"jrh3k5/autonabber/client/ynab/http"
	"jrh3k5/autonabber/client/ynab/model"
	"jrh3k5/autonabber/delta"
	"jrh3k5/autonabber/input"
	"os"

	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

// TODO: add validation of input in model.go
// TODO: warn if any of the inputs are not present in the budget to which they are being applied
// TODO: add support for parseable average-spent-##m

func main() {
	_l, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("unable to instantiate logger: %v\n", err)
		os.Exit(1)
	}
	logger := _l.Sugar()

	appArgs, err := args.GetArgs()
	if err != nil {
		logger.Fatalf("Unable to parse application arguments: %v", err)
	}

	client, err := http.NewClient(appArgs.AccessToken)
	if err != nil {
		logger.Fatalf("Unable to instantiate YNAB client: %v", err)
	}

	budget, err := getBudget(client)
	if err != nil {
		logger.Fatalf("Unable to successfully choose a budget: %v", err)
	}

	budgetCategoryGroups, err := client.GetCategories(budget)
	if err != nil {
		logger.Fatalf("Unable to retrieve budget categories: %w", err)
	}

	budgetChange, err := getBudgetChanges(appArgs.InputFile)
	if err != nil {
		logger.Fatalf("Failed to select budget change: %w", err)
	}

	ignoreMismatches, err := warnApplicationMismatches(budgetCategoryGroups, budgetChange.CategoryGroups)
	if err != nil {
		logger.Fatalf("Failed to check for and confirm acceptance of mismatches between input file and budget: %w", err)
	}

	if !ignoreMismatches {
		fmt.Println("Application has been cancelled")
		os.Exit(0)
	}

	deltas, err := delta.NewDeltas(budgetCategoryGroups, budgetChange)
	if err != nil {
		logger.Fatalf("Failed to generate delta: %w", err)
	}

	appConfirmed, err := confirmApplication(deltas)
	if err != nil {
		logger.Fatalf("Failed to get confirmation to apply changes: %w", err)
	}

	if appConfirmed {
		doAssignment, err := checkAssignability(budgetCategoryGroups, deltas)
		if err != nil {
			logger.Fatalf("Failed to check availability of assignable funds: %w", err)
		}

		if !doAssignment {
			fmt.Println("Application has been cancelled")
		} else {
			var nonZeroChanges []*delta.BudgetCategoryDelta

			for _, delta := range deltas {
				for _, change := range delta.CategoryDeltas {
					if change.HasChanges() {
						nonZeroChanges = append(nonZeroChanges, change)
						if err := client.SetBudget(budget, change.BudgetCategory, change.FinalDollars, change.FinalCents); err != nil {
							logger.Fatalf("Failed to set budget category '%s' under budget '%s' to $%d.%02d: %w", change.BudgetCategory.Name, budget.Name, change.FinalDollars, change.FinalCents, err)
						}
					}
				}
			}

			deltaDollars, deltaCents := delta.SumChanges(nonZeroChanges)
			fmt.Printf("Added $%d.%02d across %d categories\n", deltaDollars, deltaCents, len(nonZeroChanges))
		}
	} else {
		fmt.Println("Application has been cancelled")
	}

	os.Exit(0)
}

// checkAssignability checks to see if the total of the changes to be applied exceeds the amount available for assignment
// It returns true if either the user has chosen to continue or there is enough to be assigned; false if not and the application should be canceled
func checkAssignability(groups []*model.BudgetCategoryGroup, deltaGroups []*delta.BudgetCategoryDeltaGroup) (bool, error) {
	assignableDollars, assignableCents := model.GetReadyToAssign(groups)

	var deltas []*delta.BudgetCategoryDelta
	for _, group := range deltaGroups {
		for _, category := range group.CategoryDeltas {
			deltas = append(deltas, category)
		}
	}
	changeDollars, changeCents := delta.SumChanges(deltas)

	assignableTotal := assignableDollars*100 + int64(assignableCents)
	changeTotal := changeDollars*100 + int64(changeCents)

	if changeTotal > assignableTotal {
		confirmPrompt := &promptui.Prompt{
			Label:    fmt.Sprintf("Your total to be assigned ($%d.%02d) is greater than your amount ready for assignment ($%d.%02d). Do you wish to continue the application? (yes/no)", changeDollars, changeCents, assignableDollars, assignableCents),
			Validate: validateYesNo,
		}

		promptResult, err := confirmPrompt.Run()
		if err != nil {
			return false, fmt.Errorf("failed to prompt for confirmation of application: %w", err)
		}

		return promptResult == "yes", nil
	}

	return true, nil
}

func validateYesNo(input string) error {
	if input != "yes" && input != "no" {
		return fmt.Errorf("invalid selection: %s", input)
	}

	return nil
}

func confirmApplication(deltas []*delta.BudgetCategoryDeltaGroup) (bool, error) {
	delta.PrintDeltas(deltas)

	confirmPrompt := promptui.Prompt{
		Label:    "Do you wish to apply these changes? (yes/no)",
		Validate: validateYesNo,
	}

	result, err := confirmPrompt.Run()
	if err != nil {
		return false, fmt.Errorf("failed to confirm desire to apply deltas: %w", err)
	}

	return result == "yes", nil
}

func getBudget(client ynab.Client) (*model.Budget, error) {
	budgets, err := client.GetBudgets()
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}

	if len(budgets) == 0 {
		return nil, errors.New("no budgets found; please create a budget before using this tool")
	}

	if len(budgets) == 1 {
		return budgets[0], nil
	}

	budgetPromptTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "🔨 {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "✔ {{ .Name }}",
	}

	prompt := promptui.Select{
		Label:     "Select a budget",
		Templates: budgetPromptTemplate,
		Items:     budgets,
	}

	chosenBudget, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed in prompt for budget selection: %w", err)
	}

	return budgets[chosenBudget], nil
}

func getBudgetChanges(filePath string) (*input.BudgetChange, error) {
	budgetChanges, err := input.ParseInputFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to parse input file '%s': %w", filePath, err)
	}

	if len(budgetChanges.Changes) == 0 {
		return nil, errors.New("at least one change set must be supplied")
	}

	if len(budgetChanges.Changes) == 1 {
		return budgetChanges.Changes[0], nil
	}

	changePromptTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🔨 {{ .Name | cyan }}",
		Inactive: "  {{ .Name }}",
		Selected: "✔ {{ .Name }}",
	}

	prompt := promptui.Select{
		Label:     "Select a budget change",
		Templates: changePromptTemplate,
		Items:     budgetChanges.Changes,
	}

	chosenBudgetChange, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed in prompt for budget change selection: %w", err)
	}

	return budgetChanges.Changes[chosenBudgetChange], nil
}

func warnApplicationMismatches(budgetCategoryGroups []*model.BudgetCategoryGroup, inputGroups []*input.BudgetCategoryGroup) (bool, error) {
	var missingGroups []string
	missingCategoriesByGroupName := make(map[string][]string)

	for _, inputCategory := range inputGroups {
		var matchingBudgetGroup *model.BudgetCategoryGroup
		for _, budgetGroup := range budgetCategoryGroups {
			if budgetGroup.Name == inputCategory.Name {
				matchingBudgetGroup = budgetGroup
				break
			}
		}

		if matchingBudgetGroup == nil {
			missingGroups = append(missingGroups, inputCategory.Name)
			continue
		}

		budgetCategoriesByName := make(map[string]*model.BudgetCategory)
		for _, budgetCategory := range matchingBudgetGroup.Categories {
			budgetCategoriesByName[budgetCategory.Name] = budgetCategory
		}

		for _, inputCategory := range inputCategory.Changes {
			if _, categoryExists := budgetCategoriesByName[inputCategory.Name]; !categoryExists {
				var missingCategories []string
				if existingCategories, ok := missingCategoriesByGroupName[matchingBudgetGroup.Name]; ok {
					missingCategories = existingCategories
				}
				missingCategories = append(missingCategories, inputCategory.Name)
				missingCategoriesByGroupName[matchingBudgetGroup.Name] = missingCategories
			}
		}
	}

	// if there are no mistmatches, we can continue on
	if len(missingGroups) == 0 && len(missingCategoriesByGroupName) == 0 {
		return true, nil
	}

	fmt.Printf("WARNING: %d categories were in the given file that do not exist and/or %d categories specified in the input file do not exist in the budget:\n", len(missingGroups), len(missingCategoriesByGroupName))
	for _, missingGroup := range missingGroups {
		fmt.Printf("  Missing category group: %s\n", missingGroup)
	}

	for categoryGroupName, categories := range missingCategoriesByGroupName {
		fmt.Printf("  Category group with missing category/categories: %s\n", categoryGroupName)
		for _, category := range categories {
			fmt.Printf("    Missing category: %s\n", category)
		}
	}

	fmt.Println("None of the changes specified in the above category groups and categories will be applied to the budget.")

	continuePrompt := promptui.Prompt{
		Label:    "Do you wish to continue? (yes/no)",
		Validate: validateYesNo,
	}

	promptResult, err := continuePrompt.Run()
	if err != nil {
		return false, fmt.Errorf("failed to prompt user to confirm continuing with missing category groups and/or categories: %w", err)
	}

	return promptResult == "yes", nil
}
