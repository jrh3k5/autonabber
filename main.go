package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jrh3k5/autonabber/args"
	"github.com/jrh3k5/autonabber/client/ynab"
	"github.com/jrh3k5/autonabber/client/ynab/http"
	"github.com/jrh3k5/autonabber/client/ynab/model"
	"github.com/jrh3k5/autonabber/delta"
	"github.com/jrh3k5/autonabber/format"
	"github.com/jrh3k5/autonabber/input"
	"github.com/jrh3k5/oauth-cli/pkg/auth"

	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

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

	oauthToken, err := auth.DefaultGetOAuthToken(ctx,
		"https://app.ynab.com/oauth/authorize",
		"https://api.ynab.com/oauth/token",
		auth.WithLogger(logger.Infof),
		auth.WithOAuthServerPort(appArgs.OAuthServerPort))
	if err != nil {
		logger.Fatalf("Unable to retrieve OAuth token: %v", err)
	}

	var client ynab.Client
	client, err = http.NewClient(oauthToken.AccessToken)
	if err != nil {
		logger.Fatalf("Unable to instantiate YNAB client: %v", err)
	}

	if appArgs.DryRun {
		logger.Info("Dry run mode is enabled; no changes will be written to YNAB")
		client = ynab.NewReadOnlyClient(logger, client)
	}

	budget, err := getBudget(client)
	if err != nil {
		logger.Fatalf("Unable to successfully choose a budget: %v", err)
	}

	budgetCategoryGroups, err := client.GetCategories(budget)
	if err != nil {
		logger.Fatalf("Unable to retrieve budget categories: %w", err)
	}

	if appArgs.PrintBudget {
		logger.Infof("Printing budget as requested:")
		model.PrintBudgetCategoryGroups(budgetCategoryGroups, appArgs.PrintHiddenCategories)
	}

	budgetChange, err := getBudgetChanges(appArgs.ConfigFilePath)
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

	deltas, err := delta.NewDeltas(client, budget, budgetCategoryGroups, budgetChange)
	if err != nil {
		logger.Fatalf("Failed to generate delta: %w", err)
	}

	appConfirmed, err := confirmApplication(deltas)
	if err != nil {
		logger.Fatalf("Failed to get confirmation to apply changes: %w", err)
	}

	if appConfirmed {
		assignableDollars, assignableCents, err := client.GetReadyToAssign(budget)
		if err != nil {
			logger.Fatalf("Failed to get assignable dollars and cents: %w", err)
		}

		doAssignment, err := checkAssignability(assignableDollars, assignableCents, budgetCategoryGroups, deltas)
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
					}
				}
			}

			for changeIndex, change := range nonZeroChanges {
				logger.Infof("Applying change %d of %d", changeIndex+1, len(nonZeroChanges))
				if err := client.SetBudget(budget, change.BudgetCategory, change.FinalBudgetDollars, change.FinalBudgetCents); err != nil {
					formattedFinal := format.FormatUSD(change.FinalDollars, change.FinalCents)
					logger.Fatalf("Failed to set budget category '%s' under budget '%s' to %s: %w", change.BudgetCategory.Name, budget.Name, formattedFinal, err)
				}
			}

			deltaDollars, deltaCents := delta.SumChanges(nonZeroChanges)
			formattedDelta := format.FormatUSD(deltaDollars, deltaCents)
			fmt.Printf("Added %s across %d categories\n", formattedDelta, len(nonZeroChanges))
		}
	} else {
		fmt.Println("Application has been cancelled")
	}

	os.Exit(0)
}

// checkAssignability checks to see if the total of the changes to be applied exceeds the amount available for assignment
// It returns true if either the user has chosen to continue or there is enough to be assigned; false if not and the application should be canceled
func checkAssignability(assignableDollars int64, assignableCents int16, groups []*model.BudgetCategoryGroup, deltaGroups []*delta.BudgetCategoryDeltaGroup) (bool, error) {
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
		formattedChange := format.FormatUSD(changeDollars, changeCents)
		formattedAssignable := format.FormatUSD(assignableDollars, assignableCents)
		confirmPrompt := &promptui.Prompt{
			Label:    fmt.Sprintf("Your total to be assigned (%s) is greater than your amount ready for assignment (%s). Do you wish to continue the application? (yes/no)", formattedChange, formattedAssignable),
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
