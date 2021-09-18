package printer_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/mimuret/golang-iij-dpf/pkg/types"
)

var _ = Describe("core", func() {
	var (
		p       printer.HumanReadablePrinter
		headers []interface{}
		row     []interface{}
	)
	Context("Contract", func() {
		var (
			s = &core.Contract{
				Id:          "hogehoge",
				ServiceCode: "dpf00001",
				State:       types.StateBeforeStart,
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ContractId", "ServiceCode", "State"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"hogehoge", "dpf00001", "BeforeStart"}))
		})
	})
	Context("ContractList", func() {
		var (
			s = &core.ContractList{
				Items: []core.Contract{
					{
						Id:          "hogehoge",
						ServiceCode: "dpf00001",
						State:       types.StateRunning,
					},
				},
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s.Items[0])
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ContractId", "ServiceCode", "State"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"hogehoge", "dpf00001", "Started"}))
		})
	})
	Context("DelegationList", func() {
		var (
			atTime, _ = types.ParseTime(time.RFC3339Nano, "2021-06-20T07:55:17.753Z")
			s         = &core.DelegationList{
				Items: []core.Delegation{
					{
						Id:                    "m1",
						ServiceCode:           "dpm000001",
						Name:                  "example.jp.",
						Network:               "",
						DelegationRequestedAt: atTime,
					},
				},
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s.Items[0])
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ZoneId", "ServiceCode", "Name", "Network", "LastRequestTime"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"m1", "dpm000001", "example.jp.", "", "2021-06-20 07:55:17.753 +0000 UTC"}))
		})
	})
	Context("Job", func() {
		var (
			s = &core.Job{
				RequestId:    "1ADC86BA65404664B8080904E88CF7B0",
				Status:       core.JobStatusFailed,
				ErrorType:    "hoge",
				ErrorMessage: "error message",
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"RequestId", "Status", "ErrorType", "ErrorMessage"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"1ADC86BA65404664B8080904E88CF7B0", "FAILED", "hoge", "error message"}))
		})
	})
	Context("Zone", func() {
		var (
			s = &core.Zone{
				Id:               "m1",
				CommonConfigId:   1,
				ServiceCode:      "dpm0000001",
				State:            types.StateBeforeStart,
				Favorite:         types.FavoriteHighPriority,
				Name:             "example.jp.",
				Network:          "",
				Description:      "zone 1",
				ZoneProxyEnabled: types.Disabled,
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ZoneId", "ServiceCode", "Name", "State", "CommonConfigId"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"m1", "dpm0000001", "example.jp.", "BeforeStart", "1"}))
		})
	})
	Context("ZoneList", func() {
		var (
			s = &core.ZoneList{
				Items: []core.Zone{
					{
						Id:               "m1",
						CommonConfigId:   1,
						ServiceCode:      "dpm0000001",
						State:            types.StateBeforeStart,
						Favorite:         types.FavoriteHighPriority,
						Name:             "example.jp.",
						Network:          "",
						Description:      "zone 1",
						ZoneProxyEnabled: types.Disabled,
					},
				},
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s.Items[0])
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ZoneId", "ServiceCode", "Name", "State", "CommonConfigId"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"m1", "dpm0000001", "example.jp.", "BeforeStart", "1"}))
		})
	})
})
