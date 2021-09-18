package printer_test

import (
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/common_configs"
	"github.com/mimuret/golang-iij-dpf/pkg/testtool"
	"github.com/mimuret/golang-iij-dpf/pkg/types"
)

var _ = Describe("common_configs", func() {
	var (
		p       printer.HumanReadablePrinter
		headers []interface{}
		row     []interface{}
	)
	Context("CcPrimary", func() {
		var (
			s = &common_configs.CcPrimary{
				Id:      3,
				Address: net.ParseIP("192.168.0.1"),
				TsigId:  9,
				Enabled: types.Enabled,
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"Id", "Address", "TsigId", "Enabled"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "192.168.0.1", "9", "Enabled"}))
		})
	})
	Context("CcPrimaryList", func() {
		var (
			s = &common_configs.CcPrimaryList{
				Items: []common_configs.CcPrimary{
					{
						Id:      3,
						Address: net.ParseIP("192.168.0.1"),
						TsigId:  9,
						Enabled: types.Enabled,
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
			Expect(headers).To(Equal([]interface{}{"Id", "Address", "TsigId", "Enabled"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "192.168.0.1", "9", "Enabled"}))
		})
	})
	Context("CcSecNotifiedServer", func() {
		var (
			s = &common_configs.CcSecNotifiedServer{
				Id:      3,
				Address: net.ParseIP("192.168.0.1"),
				TsigId:  9,
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"Id", "Address", "TsigId"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "192.168.0.1", "9"}))
		})
	})
	Context("CcSecNotifiedServerList", func() {
		var (
			s = &common_configs.CcSecNotifiedServerList{
				Items: []common_configs.CcSecNotifiedServer{
					{
						Id:      3,
						Address: net.ParseIP("192.168.0.1"),
						TsigId:  9,
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
			Expect(headers).To(Equal([]interface{}{"Id", "Address", "TsigId"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "192.168.0.1", "9"}))
		})
	})
	Context("CcSecTransferAcl", func() {
		var (
			s = &common_configs.CcSecTransferAcl{
				Id:      3,
				Network: testtool.MustParseIPNet("192.168.0.0/24"),
				TsigId:  9,
			}
		)
		BeforeEach(func() {
			p = printer.GetHumanReadablePrinter(s)
			headers = p.GetHeaders()
			row = p.GetRow(s)
		})
		It("returns headers", func() {
			Expect(headers).To(Equal([]interface{}{"Id", "Network", "TsigId"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "192.168.0.0/24", "9"}))
		})
	})
	Context("CcSecTransferAclList", func() {
		var (
			s = &common_configs.CcSecTransferAclList{
				Items: []common_configs.CcSecTransferAcl{
					{
						Id:      3,
						Network: testtool.MustParseIPNet("192.168.0.0/24"),
						TsigId:  9,
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
			Expect(headers).To(Equal([]interface{}{"Id", "Network", "TsigId"}))
		})
		It("returns row", func() {
			Expect(row).To(Equal([]interface{}{"3", "192.168.0.0/24", "9"}))
		})
	})
})
