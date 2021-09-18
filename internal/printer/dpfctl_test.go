package printer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
)

var _ = Describe("CommandResults", func() {
	var (
		p       printer.HumanReadablePrinter
		headers []interface{}
		row     []interface{}
		s       *utils.CommandResults
	)
	Context("CurrentRecordList", func() {
		BeforeEach(func() {
			s = &utils.CommandResults{
				Items: []utils.CommandResult{{
					RequestId: "A499DEC89409406F9150329553A9AC96",
				}, {
					RequestId: "B5D25B37EB164BCBAD0A595F867A06CC",
					Job: &core.Job{
						RequestId:    "B5D25B37EB164BCBAD0A595F867A06CC",
						Status:       core.JobStatusFailed,
						ErrorType:    "hoge",
						ErrorMessage: "error message",
					},
				}},
			}
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"RequestId", "Status", "ErrorType", "ErrorMessage"}))
		})
		When("async response", func() {
			BeforeEach(func() {
				row = p.GetRow(s.Items[0])
			})
			It("returns row", func() {
				Expect(row).To(Equal([]interface{}{"A499DEC89409406F9150329553A9AC96", "", "", ""}))
			})
		})
		When("sync response", func() {
			BeforeEach(func() {
				row = p.GetRow(s.Items[1])
			})
			It("returns row", func() {
				Expect(row).To(Equal([]interface{}{"B5D25B37EB164BCBAD0A595F867A06CC", "FAILED", "hoge", "error message"}))
			})
		})
	})
})
