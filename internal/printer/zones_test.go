package printer_test

import (
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/zones"
	"github.com/mimuret/golang-iij-dpf/pkg/types"
)

var _ = Describe("zones", func() {
	var (
		p       printer.HumanReadablePrinter
		headers []interface{}
		row     []interface{}
		s1, s2  zones.Record
	)
	BeforeEach(func() {
		s1 = zones.Record{
			AttributeMeta: zones.AttributeMeta{
				ZoneID: "m1",
			},
			ID:     "r1",
			Name:   "www.example.jp.",
			TTL:    30,
			RRType: zones.TypeA,
			RData: []zones.RecordRDATA{
				{Value: "192.168.1.1"},
				{Value: "192.168.1.2"},
			},
			State:       zones.RecordStateApplied,
			Description: "SERVER(IPv4)",
			Operator:    "user1",
		}
		s2 = zones.Record{
			AttributeMeta: zones.AttributeMeta{
				ZoneID: "m1",
			},
			ID:     "r2",
			Name:   "www.example.jp.",
			TTL:    30,
			RRType: zones.TypeAAAA,
			RData: []zones.RecordRDATA{
				{Value: "2001:db8::1"},
				{Value: "2001:db8::2"},
			},
			State:       zones.RecordStateToBeAdded,
			Description: "SERVER(IPv6)",
			Operator:    "user1",
		}
	})
	Context("CurrentRecordList", func() {
		BeforeEach(func() {
			s := &zones.CurrentRecordList{
				Items: []zones.Record{s1, s2},
			}
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s.Items[0])
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ID", "Name", "TTL", "RRtype", "RData"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"r1", "www.example.jp.", "30", "A", "192.168.1.1,192.168.1.2"}))
		})
	})
	Context("Dnssec", func() {
		BeforeEach(func() {
			s := &zones.Dnssec{
				AttributeMeta: zones.AttributeMeta{
					ZoneID: "m1",
				},
				Enabled: types.Enabled,
				State:   zones.DnssecStateEnabling,
				DsState: zones.DSStateDisclose,
			}
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"Enabled", "State", "DsState"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"Enabled", "Enabling", "Disclose"}))
		})
	})
	Context("DsRecordList", func() {
		var (
			atTime, _ = types.ParseTime(time.RFC3339Nano, "2021-06-20T10:23:51.071Z")
			s         = &zones.DsRecordList{
				Items: []zones.DsRecord{
					{
						RRSet:     "46369 8 2 39F054DCB3EC1E93D8AE6D8F1AAAD91794055EA36895045FAF6F65F0 2FEBC579",
						TransitAt: atTime,
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
			Expect(headers).To(Equal([]interface{}{"TransitAt", "RDATA"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"2021-06-20 10:23:51.071 +0000 UTC", "46369 8 2 39F054DCB3EC1E93D8AE6D8F1AAAD91794055EA36895045FAF6F65F0 2FEBC579"}))
		})
	})
	Context("ManagedDnsList", func() {
		BeforeEach(func() {
			s := &zones.ManagedDnsList{
				Items: []string{
					"ns1.example.jp.",
					"ns1.example.net.",
					"ns1.example.com.",
				},
			}
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s.Items[0])
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ServerName"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"ns1.example.jp."}))
		})
	})
	Context("Record", func() {
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(&s1)
			headers = p.GetHeaders()
			row = p.GetRow(&s1)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ID", "Name", "TTL", "RRtype", "RData", "State", "Operator"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"r1", "www.example.jp.", "30", "A", "192.168.1.1,192.168.1.2", "Applied", "user1"}))
		})
	})
	Context("RecordList", func() {
		BeforeEach(func() {
			s := &zones.RecordList{
				Items: []zones.Record{s1, s2},
			}
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s.Items[0])
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ID", "Name", "TTL", "RRtype", "RData", "State", "Operator"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"r1", "www.example.jp.", "30", "A", "192.168.1.1,192.168.1.2", "Applied", "user1"}))
		})
	})
	Context("DefaultTTL", func() {
		var (
			s = &zones.DefaultTTL{
				AttributeMeta: zones.AttributeMeta{
					ZoneID: "m1",
				},
				Value:    3600,
				State:    zones.DefaultTTLStateApplied,
				Operator: "user1",
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"Value", "State", "Operator"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3600", "Applied", "user1"}))
		})
	})
	Context("HistoryList", func() {
		var (
			atTime, _ = types.ParseTime(time.RFC3339Nano, "2021-06-20T10:23:51.071Z")
			s         = &zones.HistoryList{
				Items: []zones.History{
					{
						ID:          1,
						CommittedAt: atTime,
						Description: "commit 1",
						Operator:    "user1",
					},
					{
						ID:          2,
						CommittedAt: atTime,
						Description: "commit 2",
						Operator:    "user2",
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
			Expect(headers).To(Equal([]interface{}{"ID", "CommittedAt", "Operator", "Description"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"1", "2021-06-20 10:23:51.071 +0000 UTC", "user1", "commit 1"}))
		})
	})
	Context("LogList", func() {
		var (
			atTime, _ = types.ParseTime(time.RFC3339Nano, "2021-06-20T07:55:17.753Z")
			s         = &zones.LogList{
				Items: []core.Log{
					{
						Time:      atTime,
						LogType:   "service",
						Operator:  "user1",
						Operation: "add_cc_primary",
						Target:    "1",
						Status:    core.LogStatusStart,
						RequestID: "C694DCE2D20F46E5A8DCE9EA43042B06",
					}, {
						Time:      atTime,
						LogType:   "common_config",
						Operator:  "user2",
						Operation: "create_tsig",
						Target:    "2",
						Status:    core.LogStatusSuccess,
						RequestID: "383C5E4F7968420AAE67A1636CF80497",
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
			Expect(headers).To(Equal([]interface{}{"RequestID", "Time", "Status"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"C694DCE2D20F46E5A8DCE9EA43042B06", "2021-06-20 07:55:17.753 +0000 UTC", "start"}))
		})
	})
	Context("ZoneProxy", func() {
		var (
			s = &zones.ZoneProxy{
				AttributeMeta: zones.AttributeMeta{
					ZoneID: "m1",
				},
				Enabled: types.Enabled,
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"Enabled"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"Enabled"}))
		})
	})
	Context("ZoneProxyHealthCheckList", func() {
		var (
			s = &zones.ZoneProxyHealthCheckList{
				Items: []zones.ZoneProxyHealthCheck{
					{
						Address:  net.ParseIP("192.168.0.1"),
						Status:   zones.ZoneProxyStatusSuccess,
						TsigName: "hoge",
						Enabled:  types.Disabled,
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
			Expect(headers).To(Equal([]interface{}{"Address", "Status", "TsigName", "Enabled"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"192.168.0.1", "success", "hoge", "Disabled"}))
		})
	})
	Context("Contract", func() {
		var (
			s = &zones.Contract{
				Contract: core.Contract{
					ID:          "hogehoge",
					ServiceCode: "dpf00001",
					State:       types.StateBeforeStart,
				},
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"ContractID", "ServiceCode", "State"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"hogehoge", "dpf00001", "BeforeStart"}))
		})
	})
})
