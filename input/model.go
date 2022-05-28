package input

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	addOpRegexp = regexp.MustCompile("\\+([0-9]+)(\\.([0-9]{2}))?")
)

// TODO: add validation
type BudgetChanges struct {
	Changes []*BudgetChange
}

type BudgetChange struct {
	Name           string
	CategoryGroups []*BudgetCategoryGroup
}

type BudgetCategoryGroup struct {
	Name    string
	Changes []*BudgetCategoryChange
}

func NewBudgetCategoryChange(name string, changeOperation string) *BudgetCategoryChange {
	return &BudgetCategoryChange{
		Name:            name,
		changeOperation: changeOperation,
	}
}

type BudgetCategoryChange struct {
	Name            string
	changeOperation string
}

// ApplyDelta will apply the change described in the budget category change
// It returns the given dollars and cents after this change has been applied
func (bcc *BudgetCategoryChange) ApplyDelta(dollars int64, cents int16) (int64, int16) {
	// For now, this only cares about addition, so that's easy!
	if addOpRegexp.MatchString(bcc.changeOperation) {
		regexMatch := addOpRegexp.FindStringSubmatch(bcc.changeOperation)
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

		return finalDollars, int16(finalCents)
	}

	// Not graceful, but prior validation should prevent us from reaching this point
	panic(fmt.Sprintf("Unexpected operation in change for category '%s': %s", bcc.Name, bcc.changeOperation))
}
