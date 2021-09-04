package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/utils"
)

var _ = Describe("slice", func() {
	Context("StringSliceToInterfaceSlice", func() {
		var (
			src []string
			res []interface{}
		)
		BeforeEach(func() {
			src = []string{"A", "B", "1", ""}
			res = utils.StringSliceToInterfaceSlice(src)
		})
		It("returns []interface{}", func() {
			Expect(res).To(Equal([]interface{}{"A", "B", "1", ""}))
		})
	})
})
