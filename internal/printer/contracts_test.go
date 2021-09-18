package printer_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/contracts"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/mimuret/golang-iij-dpf/pkg/types"
)

var _ = Describe("contracts", func() {
	var (
		p       printer.HumanReadablePrinter
		headers []interface{}
		row     []interface{}
	)
	Context("CommonConfig", func() {
		var (
			s = &contracts.CommonConfig{
				Id:                3,
				Name:              "共通設定1",
				ManagedDNSEnabled: types.Disabled,
				Default:           types.Enabled,
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"Id", "Name", "ManagedDNSEnabled", "Default"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "共通設定1", "Disabled", "Enabled"}))
		})
	})
	Context("CommonConfigList", func() {
		var (
			s = &contracts.CommonConfigList{
				Items: []contracts.CommonConfig{
					{
						Id:                3,
						Name:              "共通設定1",
						ManagedDNSEnabled: types.Disabled,
						Default:           types.Enabled,
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
			Expect(headers).To(Equal([]interface{}{"Id", "Name", "ManagedDNSEnabled", "Default"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "共通設定1", "Disabled", "Enabled"}))
		})
	})
	Context("LogList", func() {
		var (
			atTime, _ = types.ParseTime(time.RFC3339Nano, "2021-06-20T07:55:17.753Z")
			s         = &contracts.LogList{
				Items: []core.Log{
					{
						Time:      atTime,
						LogType:   "service",
						Operator:  "user1",
						Operation: "add_cc_primary",
						Target:    "1",
						Status:    core.LogStatusStart,
						RequestId: "C694DCE2D20F46E5A8DCE9EA43042B06",
					}, {
						Time:      atTime,
						LogType:   "common_config",
						Operator:  "user2",
						Operation: "create_tsig",
						Target:    "2",
						Status:    core.LogStatusSuccess,
						RequestId: "383C5E4F7968420AAE67A1636CF80497",
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
			Expect(headers).To(Equal([]interface{}{"RequestId", "Time", "Status"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"C694DCE2D20F46E5A8DCE9EA43042B06", "2021-06-20 07:55:17.753 +0000 UTC", "start"}))
		})
	})
	Context("ContractPartnerList", func() {
		var (
			s = &contracts.ContractPartnerList{
				Items: []contracts.ContractPartner{
					{
						ServiceCode: "svc1",
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
			Expect(headers).To(Equal([]interface{}{"ServiceCode"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"svc1"}))
		})
	})
	Context("QpsHistoryList", func() {
		var (
			s = &contracts.QpsHistoryList{
				Items: []contracts.QpsHistory{
					{
						ServiceCode: "svc1",
						Name:        "example.jp.",
						Values: []contracts.QpsValue{
							{
								Month: "202108",
								Qps:   100,
							},
							{
								Month: "202109",
								Qps:   200,
							},
						},
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
			Expect(headers).To(Equal([]interface{}{"ServiceCode", "Name", "LastMonth", "LastQps"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"svc1", "example.jp.", "202109", "200"}))
		})
	})
	Context("ContractZoneList", func() {
		var (
			s = &contracts.ContractZoneList{
				Items: []core.Zone{
					{
						Id:             "m1",
						ServiceCode:    "dpm000001",
						Name:           "example.jp.",
						State:          types.StateRunning,
						CommonConfigId: 1,
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
			Expect(row).To(Equal([]interface{}{"m1", "dpm000001", "example.jp.", "Started", "1"}))
		})
	})
})
