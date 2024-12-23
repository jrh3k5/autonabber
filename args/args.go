package args

import (
	"flag"

	"github.com/jrh3k5/autonabber/client/auth"
)

type Args struct {
	PrintBudget           bool
	PrintHiddenCategories bool
	DryRun                bool
	OAuthServerPort       int
}

func GetArgs() (*Args, error) {
	var printBudget bool
	flag.BoolVar(&printBudget, "print-budget", false, "whether or not the budget should be printed to the screen")

	var printHiddenCategories bool
	flag.BoolVar(&printHiddenCategories, "print-hidden-categories", false, "if print-budget is specified, controls if hidden categories are printed")

	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "changes to budgets are not actually applied")

	var oAuthServerPort int
	flag.IntVar(&oAuthServerPort, "oauth-server-port", auth.DefaultOAuthServerPort, "the port on which the OAuth server should run")

	flag.Parse()

	return &Args{
		PrintBudget:           printBudget,
		PrintHiddenCategories: printHiddenCategories,
		DryRun:                dryRun,
	}, nil
}
