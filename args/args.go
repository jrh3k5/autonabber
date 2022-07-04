package args

import (
	"errors"
	"flag"
)

type Args struct {
	AccessToken           string
	InputFile             string
	PrintBudget           bool
	PrintHiddenCategories bool
}

func GetArgs() (*Args, error) {
	var accessToken string
	flag.StringVar(&accessToken, "access-token", "", "the access token to be used to communication with YNAB")

	var inputFile string
	flag.StringVar(&inputFile, "file", "", "the file containing the changes to be applied")

	var printBudget bool
	flag.BoolVar(&printBudget, "print-budget", false, "whether or not the budget should be printed to the screen")

	var printHiddenCategories bool
	flag.BoolVar(&printHiddenCategories, "print-hidden-categories", false, "if print-budget is specified, controls if hidden categories are printed")

	flag.Parse()

	if accessToken == "" {
		return nil, errors.New("an access token must be provided")
	}

	if inputFile == "" {
		return nil, errors.New("a file containing budget changes must be provided")
	}

	return &Args{
		AccessToken:           accessToken,
		InputFile:             inputFile,
		PrintBudget:           printBudget,
		PrintHiddenCategories: printHiddenCategories,
	}, nil
}
