package utils_test

import (
	_ "github.com/mimuret/dpfctl/internal/params"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

var _ = Describe("command", func() {
	Context("ValidArgsFunction", func() {
		var (
			f []string
			d cobra.ShellCompDirective
		)
		When("is support action", func() {
			BeforeEach(func() {
				f, d = utils.ValidArgsFunction(api.ActionRead, []string{})
			})
			It("returns api names", func() {
				Expect(len(f)).NotTo(BeZero())
			})
			It("return ShellCompDirectiveDefault", func() {
				Expect(d).To(Equal(cobra.ShellCompDirectiveDefault))
			})
		})
		When("is not support action", func() {
			BeforeEach(func() {
				f, d = utils.ValidArgsFunction(api.ActionCount, []string{"hoge"})
			})
			It("returns api names", func() {
				Expect(len(f)).To(BeZero())
			})
			It("return ShellCompDirectiveNoSpace", func() {
				Expect(d).To(Equal(cobra.ShellCompDirectiveNoSpace))
			})
		})
	})
})
