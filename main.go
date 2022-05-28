package main

import (
	"fmt"
	"jrh3k5/autonabber/args"
	"jrh3k5/autonabber/client/ynab/http"
	"os"

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

	budgets, err := client.GetBudgets()
	if err != nil {
		logger.Errorf("Unable to retrieve budgets: %v", err)
		os.Exit(1)
	}

	for _, budget := range budgets {
		logger.Info("Budget found: " + budget.Name)
	}

	os.Exit(0)
}
