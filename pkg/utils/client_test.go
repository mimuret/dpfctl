package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

var _ = Describe("client", func() {
	Context("NewClient", func() {
		var (
			cl  api.ClientInterface
			err error
		)
		BeforeEach(func() {
			cl, err = utils.NewClient(nil)
		})
		It("returns github.com/mimuret/golang-iij-dpf/pkg/api.Client", func() {
			Expect(cl).NotTo(BeNil())
		})
		It("not return err", func() {
			Expect(err).NotTo(BeNil())
		})
	})
})
