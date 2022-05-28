package main

import (
	"errors"
	"fmt"
	"jrh3k5/autonabber/args"
	"jrh3k5/autonabber/client/ynab"
	"jrh3k5/autonabber/client/ynab/http"
	"jrh3k5/autonabber/client/ynab/model"
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
		logger.Errorf("Unable to parse application arguments: %v", err)
		os.Exit(1)
	}

	client, err := http.NewClient(appArgs.AccessToken)
	if err != nil {
		logger.Errorf("Unable to instantiate YNAB client: %v", err)
		os.Exit(1)
	}

	budget, err := getBudget(client)
	if err != nil {
		logger.Errorf("Unable to successfully choose a budget: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Chosen budget: %s\n", budget.Name)

	os.Exit(0)
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
