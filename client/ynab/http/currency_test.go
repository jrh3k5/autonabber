package http_test

import (
	"github.com/jrh3k5/autonabber/client/ynab/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Currency", func() {
	Context("ParseMillidollars", func() {
		It("correctly parses positive dollar amounts", func() {
			dollars, cents := http.ParseMillidollars(123930)
			Expect(dollars).To(Equal(int64(123)), "dollar should match")
			Expect(cents).To(Equal(int16(93)), "cent amount should match")
		})

		It("correctly parses negative dollar amounts", func() {
			dollars, cents := http.ParseMillidollars(-123930)
			Expect(dollars).To(Equal(int64(-123)), "dollar should match")
			Expect(cents).To(Equal(int16(-93)), "cent amount should match")
		})
	})

	Context("ToMilliDollars", func() {
		It("correctly converts positive dollar amounts", func() {
			Expect(http.ToMillidollars(123, 93)).To(Equal(int64(123930)))
		})

		It("correctly converts negative dollar amounts", func() {
			Expect(http.ToMillidollars(-123, -93)).To(Equal(int64(-123930)))
		})

		It("correctly converts negative cent amounts", func() {
			Expect(http.ToMillidollars(0, -22)).To(Equal(int64(-220)))
		})
	})
})
