package utils_test

import (
	"fmt"

	_ "github.com/mimuret/dpfctl/internal/params"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/mimuret/golang-iij-dpf/pkg/testtool"
)

var _ = Describe("command", func() {
	Context("NewCommend", func() {
		var c *cobra.Command
		BeforeEach(func() {
			c = utils.NewCommand("get", api.ActionRead, func(cmd *cobra.Command, cl api.ClientInterface, args []string, resources []apis.Spec) error {
				return nil
			})
		})
		It("set usage", func() {
			Expect(c.Use).To(Equal("get"))
		})
		Context("RunE", func() {
			When("set filename", func() {
				When("empty file", func() {
					JustBeforeEach(func() {
						viper.Set("filename", "empty.yaml")
					})
					JustAfterEach(func() {
						viper.Set("filename", "")
						viper.Set("filename", "testdata/single-doc.yaml")
					})
					It("returns error", func() {
						Expect(c.RunE(c, []string{})).To(HaveOccurred())
					})
				})
				When("valid file", func() {
					JustBeforeEach(func() {
						viper.Set("filename", "testdata/single-doc.yaml")
					})
					JustAfterEach(func() {
						viper.Set("filename", "")
					})
					It("returns error", func() {
						Expect(c.RunE(c, []string{})).To(Succeed())
					})
				})
			})
			When("not set subcmd", func() {
				It("returns error", func() {
					Expect(c.RunE(c, []string{})).To(HaveOccurred())
				})
			})
			When("set subcmd", func() {
				When("no set resource", func() {
					It("returns error", func() {
						Expect(c.RunE(c, []string{"zones"})).To(HaveOccurred())
					})
				})
				When("set resource", func() {
					It("returns error", func() {
						Expect(c.RunE(c, []string{"zones", "mxxxxxx"})).To(Succeed())
					})
				})
			})
			When("failed to create client", func() {
				JustBeforeEach(func() {
					utils.NewClient = func(logger api.Logger) (api.ClientInterface, error) {
						return testtool.NewTestClient("token", "http://localhost", logger), fmt.Errorf("error")
					}
				})
				JustAfterEach(func() {
					utils.NewClient = func(logger api.Logger) (api.ClientInterface, error) {
						return testtool.NewTestClient("token", "http://localhost", logger), nil
					}
				})
				It("returns error", func() {
					Expect(c.RunE(c, []string{"zones", "mxxxxxx"})).To(HaveOccurred())
				})
			})
		})
	})
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
