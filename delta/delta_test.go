package delta_test

import (
	"github.com/jrh3k5/autonabber/client/mock_ynab"
	"github.com/jrh3k5/autonabber/client/ynab/model"
	"github.com/jrh3k5/autonabber/delta"
	"github.com/jrh3k5/autonabber/input"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delta", func() {
	var mockController *gomock.Controller
	var ynabClient *mock_ynab.MockClient

	BeforeEach(func() {
		mockController = gomock.NewController(GinkgoT())
		ynabClient = mock_ynab.NewMockClient(mockController)
	})

	AfterEach(func() {
		if mockController != nil {
			mockController.Finish()
		}
	})

	It("should generate a delta", func() {
		actual := []*model.BudgetCategoryGroup{
			{
				Name: "Frequent",
				Categories: []*model.BudgetCategory{
					{
						Name:             "Groceries",
						AvailableDollars: 15,
						AvailableCents:   89,
						BudgetedDollars:  75,
						BudgetedCents:    24,
					},
					{
						Name:             "Movies",
						AvailableDollars: 24,
						AvailableCents:   16,
						BudgetedDollars:  30,
						BudgetedCents:    00,
					},
					// Hey, don't judge me for disproportionate dining out budgeting over groceries >:(
					{
						Name:             "Eating Out",
						AvailableDollars: 99,
						AvailableCents:   76,
						BudgetedDollars:  150,
						BudgetedCents:    22,
					},
				},
			},
			{

				Name: "Required",
				Categories: []*model.BudgetCategory{
					{
						Name:             "Mortgage",
						AvailableDollars: 800,
						AvailableCents:   27,
						BudgetedDollars:  10,
						BudgetedCents:    12,
					},
				},
			},
		}

		groceriesChange, err := input.NewBudgetCategoryChange("Groceries", "+14.58")
		Expect(err).To(BeNil(), "creating the Groceries change should not fail")

		eatingOutChange, err := input.NewBudgetCategoryChange("Eating Out", "+34")
		Expect(err).To(BeNil(), "creating the Eating Out change should not fail")

		toApply := &input.BudgetChange{
			CategoryGroups: []*input.BudgetCategoryGroup{
				{
					Name: "Frequent",
					Changes: []*input.BudgetCategoryChange{
						groceriesChange,
						eatingOutChange,
					},
				},
			},
		}

		delta, err := delta.NewDeltas(ynabClient, &model.Budget{}, actual, toApply)
		Expect(err).To(BeNil(), "the delta formulation should not have failed")
		Expect(delta).To(HaveLen(2), "the delta should have all of the category groups given, even if not all have changes to apply")

		// "Required" grouping
		requiredGrouping := getGroupByName("Required", delta)
		Expect(requiredGrouping).ToNot(BeNil(), "a change group for 'Required' should have been found")
		Expect(requiredGrouping.CategoryDeltas).To(HaveLen(1), "there should only be one change for the 'Required' grouping")
		mortgageDelta := getDeltaByName("Mortgage", requiredGrouping.CategoryDeltas)
		Expect(mortgageDelta).ToNot(BeNil(), "there should be a 'Required' mortgage delta")
		Expect(mortgageDelta.InitialDollars).To(Equal(int64(800)), "the initial dollars for the Mortgage should be recorded")
		Expect(mortgageDelta.InitialCents).To(Equal(int16(27)), "the initial cents for the Mortgage should be recorded")
		Expect(mortgageDelta.FinalDollars).To(Equal(mortgageDelta.InitialDollars), "because there were no changes for the mortgage, the final dollars should be the initial dollars")
		Expect(mortgageDelta.FinalCents).To(Equal(mortgageDelta.InitialCents), "because there were no changes for the mortgage, the final cents should be the initial cents")
		Expect(mortgageDelta.FinalBudgetDollars).To(Equal(int64(10)), "because there were no changes to the mortgage budget, the final budgeted dollars amount should be the initial budgeted")
		Expect(mortgageDelta.FinalBudgetCents).To(Equal(int16(12)), "because there were no changes to the mortgage budget, the final budgeted cents amount should be the initial budgeted")
		mortgageDollars, mortgageCents := mortgageDelta.CalculateDelta()
		Expect(mortgageDollars).To(Equal(int64(0)), "there should be no dollar change for the mortgage")
		Expect(mortgageCents).To(Equal(int16(0)), "there should be no cent change for the mortgage")

		// "Frequent" grouping
		frequentGrouping := getGroupByName("Frequent", delta)
		Expect(frequentGrouping).ToNot(BeNil(), "a change group for 'Frequent' should have been found")
		Expect(frequentGrouping.CategoryDeltas).To(HaveLen(3), "all three categories should be returned, even if they do not all have changes")

		groceriesDelta := getDeltaByName("Groceries", frequentGrouping.CategoryDeltas)
		Expect(groceriesDelta).ToNot(BeNil(), "there should be a Groceries delta")
		Expect(groceriesDelta.InitialDollars).To(Equal(int64(15)), "the Groceries' initial dollars should be recorded")
		Expect(groceriesDelta.InitialCents).To(Equal(int16(89)), "the Groceries' initial cents should be recorded")
		Expect(groceriesDelta.FinalDollars).To(Equal(int64(30)), "the delta should be applied to the groceries' final dollars")
		Expect(groceriesDelta.FinalCents).To(Equal(int16(47)), "the delta should be applied to the groceries' final cents")
		Expect(groceriesDelta.FinalBudgetDollars).To(Equal(int64(89)), "the budgeted groceries dollars should be initial budget + delta")
		Expect(groceriesDelta.FinalBudgetCents).To(Equal(int16(82)), "the budgeted groceries cents should initial budget + delta")
		groceriesDollars, groceriesCents := groceriesDelta.CalculateDelta()
		Expect(groceriesDollars).To(Equal(int64(14)), "the correct dollars delta for the Groceries category should be returned")
		Expect(groceriesCents).To(Equal(int16(58)), "the correct cents delta for the Groceries category should be returned")

		moviesDelta := getDeltaByName("Movies", frequentGrouping.CategoryDeltas)
		Expect(moviesDelta).ToNot(BeNil(), "there should be a Movies delta")
		Expect(moviesDelta.InitialDollars).To(Equal(int64(24)), "the initial Movies dollar amount should be recorded")
		Expect(moviesDelta.InitialCents).To(Equal(int16(16)), "the initial Movies cent amount should be recorded")
		Expect(moviesDelta.FinalDollars).To(Equal(moviesDelta.InitialDollars), "because there were no changes for the movies budget, its final dollars should be its initial dollars")
		Expect(moviesDelta.FinalCents).To(Equal(moviesDelta.FinalCents), "because there were no changes to the movies budget, its final cents should be its initial cents")
		Expect(moviesDelta.FinalBudgetDollars).To(Equal(int64(30)), "because there were no changes to the movie budget, the final budgeted dollars amount should be the initial budgeted")
		Expect(moviesDelta.FinalBudgetCents).To(Equal(int16(00)), "because there were no changes to the movie budget, the final budgeted cents amount should be the initial budgeted")
		moviesDollars, moviesCents := moviesDelta.CalculateDelta()
		Expect(moviesDollars).To(Equal(int64(0)), "there should be 0 new dollars for Movies")
		Expect(moviesCents).To(Equal(int16(0)), "there should be 0 new cents for Movies")

		eatingOutDelta := getDeltaByName("Eating Out", frequentGrouping.CategoryDeltas)
		Expect(eatingOutDelta).ToNot(BeNil(), "there should be an Eating Out delta")
		Expect(eatingOutDelta.InitialDollars).To(Equal(int64(99)), "the initial dollars for Eating Out should be recorded")
		Expect(eatingOutDelta.InitialCents).To(Equal(int16(76)), "the initial cents for Eating Out should be recorded")
		Expect(eatingOutDelta.FinalDollars).To(Equal(int64(133)), "the eating out final dollars should reflect the change")
		Expect(eatingOutDelta.FinalCents).To(Equal(int16(76)), "the eating out final cents should reflect the budget change")
		Expect(eatingOutDelta.FinalBudgetDollars).To(Equal(int64(184)), "the eating out final budget dollars should be initial + delta")
		Expect(eatingOutDelta.FinalBudgetCents).To(Equal(int16(22)), "the eating out final budget dollars should be initial + delta")
		eatingOutDollars, eatingOutCents := eatingOutDelta.CalculateDelta()
		Expect(eatingOutDollars).To(Equal(int64(34)), "the delta for Eating Out dollars should be correct")
		Expect(eatingOutCents).To(Equal(int16(0)), "the delta for Eating Out cents should be correct")
	})

	Context("for monthly average expenditures", func() {
		It("should apply the returned average to the initial value", func() {
			budget := &model.Budget{
				ID: "test-budget-id",
			}

			budgetCategory := &model.BudgetCategory{
				ID:              "test-budget-category-ID",
				Name:            "Groceries",
				BudgetedDollars: 15,
				BudgetedCents:   89,
			}
			actual := []*model.BudgetCategoryGroup{
				{
					Name: "Frequent",
					Categories: []*model.BudgetCategory{
						budgetCategory,
					},
				},
			}

			groceriesChange, err := input.NewBudgetCategoryChange("Groceries", "+average-spent-9m")
			Expect(err).To(BeNil(), "creating the groceries change should not fail")
			toApply := &input.BudgetChange{
				CategoryGroups: []*input.BudgetCategoryGroup{
					{
						Name: "Frequent",
						Changes: []*input.BudgetCategoryChange{
							groceriesChange,
						},
					},
				},
			}

			averageDollars := int64(35)
			averageCents := int16(28)
			ynabClient.EXPECT().GetMonthlyAverageSpent(gomock.Eq(budget), gomock.Eq(budgetCategory), gomock.Eq(9)).AnyTimes().Return(averageDollars, averageCents, nil)

			delta, err := delta.NewDeltas(ynabClient, budget, actual, toApply)
			Expect(err).To(BeNil(), "the delta formulation should not have failed")
			Expect(delta).To(HaveLen(1), "the delta should have the given budget categories")

			frequentGrouping := delta[0]
			Expect(frequentGrouping.CategoryDeltas).To(HaveLen(1), "there should be a groceries delta")

			deltaDollars, deltaCents := frequentGrouping.CategoryDeltas[0].CalculateDelta()
			Expect(deltaDollars).To(Equal(averageDollars), "the dollar delta should be equal to the average dollars spent")
			Expect(deltaCents).To(Equal(averageCents), "the cents delta should be equal to the average cents spent")
		})
	})
})

func getGroupByName(name string, groups []*delta.BudgetCategoryDeltaGroup) *delta.BudgetCategoryDeltaGroup {
	for _, group := range groups {
		if name == group.Name {
			return group
		}
	}

	return nil
}

func getDeltaByName(name string, changes []*delta.BudgetCategoryDelta) *delta.BudgetCategoryDelta {
	for _, change := range changes {
		if name == change.BudgetCategory.Name {
			return change
		}
	}

	return nil
}
