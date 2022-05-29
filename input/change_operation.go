package input

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/jrh3k5/autonabber/client/ynab"
	"github.com/jrh3k5/autonabber/client/ynab/model"
)

var (
	addOpRegexp                = regexp.MustCompile("\\+([0-9]+)(\\.([0-9]{2}))?")
	averageMonthlySpentPattern = "average\\-spent\\-([1-9][0-9]*)m"
	averageMonthlySpentRegex   = regexp.MustCompile("\\+" + averageMonthlySpentPattern)
)

type changeOperationFn func(change *BudgetCategoryChange, client ynab.Client, budget *model.Budget, budgetCategory *model.BudgetCategory, dollars int64, cents int16) (int64, int16, error)

func addRawValue(change *BudgetCategoryChange, _ ynab.Client, _ *model.Budget, _ *model.BudgetCategory, dollars int64, cents int16) (int64, int16, error) {
	regexMatch := addOpRegexp.FindStringSubmatch(change.changeOperation)
	initialTotal := dollars*100 + int64(cents)
	addedDollars, _ := strconv.ParseInt(regexMatch[1], 10, 64)
	var addedCents int64
	parsedCents := regexMatch[3]
	if parsedCents != "" {
		addedCents, _ = strconv.ParseInt(parsedCents, 10, 16)
	}
	finalTotal := initialTotal + addedDollars*100 + addedCents

	finalCents := finalTotal % 100
	finalDollars := (finalTotal - finalCents) / 100

	return finalDollars, int16(finalCents), nil
}

func addMonthlyAverage(change *BudgetCategoryChange, client ynab.Client, budget *model.Budget, budgetCategory *model.BudgetCategory, dollars int64, cents int16) (int64, int16, error) {
	monthLookback, _ := strconv.ParseInt(averageMonthlySpentRegex.FindStringSubmatch(change.changeOperation)[1], 10, 16)
	averageDollars, averageCents, err := client.GetMonthlyAverageSpent(budget, budgetCategory, int(monthLookback))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to retrieve average monthly spent for category '%s': %w", budgetCategory.Name, err)
	}

	totalDollars := dollars + averageDollars
	totalCents := cents + averageCents
	totalCentidollars := totalDollars*100 + int64(totalCents)

	finalCents := totalCentidollars % 100
	finalDollars := (totalCentidollars - finalCents) / 100

	return finalDollars, int16(finalCents), nil
}
