package delta_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"go.uber.org/zap"

	"github.com/jrh3k5/autonabber/client/mock_ynab"
	"github.com/jrh3k5/autonabber/client/ynab/model"
	"github.com/jrh3k5/autonabber/delta"
)

var _ = Describe("ClientApplicator", func() {
	var setBudgets []*budgetSetting
	var applicator *delta.ClientApplicator

	BeforeEach(func() {
		setBudgets = nil

		mockController := gomock.NewController(GinkgoT())

		client := mock_ynab.NewMockClient(mockController)
		client.EXPECT().SetBudget(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(budget *model.Budget, budgetCategory *model.BudgetCategory, setDollarAmount int64, setCentAmount int16) error {
			setBudgets = append(setBudgets, &budgetSetting{
				budget:          budget,
				budgetCategory:  budgetCategory,
				setDollarAmount: setDollarAmount,
				setCentAmount:   setCentAmount,
			})
			return nil
		})

		applicator = delta.NewClientApplicator(zap.NewNop(), client)
	})

	Context("ApplyChanges", func() {
		It("should apply the given final dollars and cents", func() {
			budgetCategory0 := &model.BudgetCategory{
				ID: "category-0",
			}
			change0 := &delta.BudgetCategoryDelta{
				InitialDollars:     36,
				InitialCents:       72,
				FinalDollars:       98,
				FinalCents:         23,
				FinalBudgetDollars: 61,
				FinalBudgetCents:   51,
				BudgetCategory:     budgetCategory0,
			}

			budgetCategory1 := &model.BudgetCategory{
				ID: "category-1",
			}
			change1 := &delta.BudgetCategoryDelta{
				InitialDollars:     23,
				InitialCents:       21,
				FinalDollars:       34,
				FinalCents:         22,
				FinalBudgetDollars: 17,
				FinalBudgetCents:   31,
				BudgetCategory:     budgetCategory1,
			}

			group := &delta.BudgetCategoryDeltaGroup{
				Name:           "Multiple Changes",
				CategoryDeltas: []*delta.BudgetCategoryDelta{change0, change1},
			}

			budget := &model.Budget{
				ID:   "7652b638-3d0d-11ed-b878-0242ac120002",
				Name: "Target Budget",
			}

			changeDollars, changeCents, err := applicator.ApplyChanges(context.Background(), budget, []*delta.BudgetCategoryDeltaGroup{group})
			Expect(err).To(BeNil(), "applying the changes should not fail")
			Expect(changeDollars).To(Equal(int64(72)), "the returned dollar delta should be summed across all of the changes")
			Expect(changeCents).To(Equal(int16(52)), "the return cent delta should be summed across all of the changes")

			Expect(setBudgets).To(HaveLen(2), "two budgets should have been set")
			Expect(setBudgets).To(containSetting(&budgetSetting{
				budget:          budget,
				budgetCategory:  budgetCategory0,
				setDollarAmount: change0.FinalBudgetDollars,
				setCentAmount:   change0.FinalBudgetCents,
			}), "the budget settings should include an application of change 0")
			Expect(setBudgets).To(containSetting(&budgetSetting{
				budget:          budget,
				budgetCategory:  budgetCategory1,
				setDollarAmount: change1.FinalBudgetDollars,
				setCentAmount:   change1.FinalBudgetCents,
			}), "the budget settings should include an application of change 1")
		})

		It("should skip deltas that have no actual changes", func() {
			noChanges := &delta.BudgetCategoryDelta{
				InitialDollars:     36,
				InitialCents:       72,
				FinalDollars:       36,
				FinalCents:         72,
				FinalBudgetDollars: 36,
				FinalBudgetCents:   72,
				BudgetCategory: &model.BudgetCategory{
					ID: "no-change-budget-category",
				},
			}

			noChangesGroup := &delta.BudgetCategoryDeltaGroup{
				Name:           "No Changes",
				CategoryDeltas: []*delta.BudgetCategoryDelta{noChanges},
			}

			budget := &model.Budget{
				ID:   "abcdef",
				Name: "Target Budget",
			}

			changeDollars, changeCents, err := applicator.ApplyChanges(context.Background(), budget, []*delta.BudgetCategoryDeltaGroup{noChangesGroup})
			Expect(err).To(BeNil(), "applying the changes should not fail")
			Expect(changeDollars).To(Equal(int64(0)), "no dollar changes should have been applied")
			Expect(changeCents).To(Equal(int16(0)), "no cent changes should have been applied")

			Expect(setBudgets).To(HaveLen(0), "no budgets should have been set")
		})
	})
})

type budgetSetting struct {
	budget          *model.Budget
	budgetCategory  *model.BudgetCategory
	setDollarAmount int64
	setCentAmount   int16
}

func (bs *budgetSetting) String() string {
	return fmt.Sprintf("{ budget: '%v', budgetCategory: '%v', setDollarAmount: %d, setCentAmount: %d }", bs.budget, bs.budgetCategory, bs.setDollarAmount, bs.setCentAmount)
}

type containsSettingMatcher struct {
	expectedSetting *budgetSetting
}

func containSetting(expectedSetting *budgetSetting) types.GomegaMatcher {
	return &containsSettingMatcher{
		expectedSetting: expectedSetting,
	}
}

// FailureMessage implements types.GomegaMatcher
func (c *containsSettingMatcher) FailureMessage(actual interface{}) string {
	actualSettings := actual.([]*budgetSetting)
	actualStrings := make([]string, len(actualSettings))
	for settingIdx, actualSetting := range actualSettings {
		actualStrings[settingIdx] = actualSetting.String()
	}
	return fmt.Sprintf("expected to find budget setting in %s in %d actual: [%s]", c.expectedSetting, len(actualSettings), strings.Join(actualStrings, ", "))
}

// Match implements types.GomegaMatcher
func (c *containsSettingMatcher) Match(actual interface{}) (bool, error) {
	actualSettings, ok := actual.([]*budgetSetting)
	if !ok {
		return false, fmt.Errorf("actual must be of type []*budgetSetting: %v", actual)
	}

	for _, actualSetting := range actualSettings {
		if actualSetting.budget == nil {
			continue
		}

		if actualSetting.budget.ID != c.expectedSetting.budget.ID {
			continue
		}

		if actualSetting.budgetCategory == nil {
			continue
		}

		if actualSetting.budgetCategory.ID != c.expectedSetting.budgetCategory.ID {
			continue
		}

		if actualSetting.setDollarAmount != c.expectedSetting.setDollarAmount ||
			actualSetting.setCentAmount != c.expectedSetting.setCentAmount {
			continue
		}

		return true, nil
	}

	return false, nil
}

// NegatedFailureMessage implements types.GomegaMatcher
func (c *containsSettingMatcher) NegatedFailureMessage(actual interface{}) string {
	actualSettings := actual.([]*budgetSetting)
	actualStrings := make([]string, len(actualSettings))
	for settingIdx, actualSetting := range actualSettings {
		actualStrings[settingIdx] = actualSetting.String()
	}
	return fmt.Sprintf("did not expect to find budget setting in %s in %d actual: [%s]", c.expectedSetting, len(actualSettings), strings.Join(actualStrings, ", "))
}
