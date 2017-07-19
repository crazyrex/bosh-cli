package table_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/ui/table"
)

var _ = Describe("Headers", func() {
	Describe("KeyifyHeader", func() {
		It("converts alphanumeric to lowercase", func() {
			keyifyHeader := table.KeyifyHeader("Header1")
			Expect(keyifyHeader).To(Equal("header1"))
		})

		It("removes '(' and ')'", func() {
			keyifyHeader := table.KeyifyHeader("Header(1)")
			Expect(keyifyHeader).To(Equal("header1"))
		})

		Context("given a header that only contains non-alphanumeric and alphanumeric", func() {
			It("should non-alphanumeric to underscore", func() {
				keyifyHeader := table.KeyifyHeader("FOO!@AND#$BAR")
				Expect(keyifyHeader).To(Equal("foo_and_bar"))
			})
		})

		Context("given a header that only contains non-alphanumeric", func() {
			It("converts to underscore", func() {
				keyifyHeader := table.KeyifyHeader("!@#$")
				Expect(keyifyHeader).To(Equal("_"))
			})

			It("converts empty header to underscore", func() {
				keyifyHeader := table.KeyifyHeader("")
				Expect(keyifyHeader).To(Equal("_"))
			})
		})
	})

	Describe("SetColumnVisibility", func() {
		Context("when given a header that does not exist", func() {
			It("returns an error", func() {
				t := table.Table{
					Header: []table.Header{table.NewHeader("header1")},
				}

				err := t.SetColumnVisibility([]table.Header{table.NewHeader("non-matching-header")})
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
