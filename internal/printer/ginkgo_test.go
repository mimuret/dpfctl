package printer_test

import (
	"testing"

	_ "github.com/mimuret/dpfctl/internal/printer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "human readable print format test Suite")
}

var _ = BeforeSuite(func() {

})

var _ = BeforeEach(func() {
})

var _ = AfterSuite(func() {
})
