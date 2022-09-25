package undo_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jrh3k5/autonabber/client/ynab/model"
	"github.com/jrh3k5/autonabber/input"
	"github.com/jrh3k5/autonabber/undo"
)

var _ = Describe("Generator", func() {
	var budget *model.Budget
	var budgetChange *input.BudgetChange
	var generator *undo.GeneratorImpl

	BeforeEach(func() {
		budget = &model.Budget{
			ID:   "82e02bf0-3d13-11ed-b878-0242ac120002",
			Name: "Targeted Budget",
		}
		budgetChange = &input.BudgetChange{
			Name: "Applied Budget Change",
		}

		generator = undo.NewGeneratorImpl()
	})

	It("should generate a set of undo actions for a given set of changes", func() {

	})
})
