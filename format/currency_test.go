package format_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jrh3k5/autonabber/format"
)

var _ = Describe("Currency", func() {
	It("pads the cents out with leading zeros", func() {
		Expect(format.FormatUSD(1, 2)).To(Equal("$1.02"))
	})

	Context("the currency is negative", func() {
		It("should format it with a leading negative sign", func() {
			Expect(format.FormatUSD(-12, -23)).To(Equal("-$12.23"))
		})

		It("should format it with a leading negative sign even with zero dollars", func() {
			Expect(format.FormatUSD(0, -23)).To(Equal("-$0.23"))
		})
	})
})
