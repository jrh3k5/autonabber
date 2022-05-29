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

	deltas, err := delta.NewDeltas(budgetCategoryGroups, budgetChange)
	if err != nil {
		logger.Fatalf("Failed to generate delta: %w", err)
	}

	doAssignment, err := checkAssignability(budgetCategoryGroups, deltas)
	if err != nil {
		logger.Fatalf("Failed to check availability of assignable funds: %w", err)
	}

	if !doAssignment {
		fmt.Println("Application has been cancelled")
	}

	appConfirmed, err := confirmApplication(deltas)
	if err != nil {
		logger.Fatalf("Failed to get confirmation to apply changes: %w", err)
	}

	if appConfirmed {
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
			Label:    fmt.Sprintf("Your total to be assigned ($%d.%02d) is greater than your amount ready for assignment ($%d.%02d). Do you wish to continue the application?", changeDollars, changeCents, assignableDollars, assignableCents),
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
