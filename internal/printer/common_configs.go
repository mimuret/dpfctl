package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/common_configs"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&common_configs.CcPrimary{}, &common_configs.CcPrimaryList{}},
		[]string{"ID", "Address", "TsigID", "Enabled"},
		[]string{"{{ .ID }}", "{{ .Address }}", "{{ .TsigID }}", "{{ .Enabled }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&common_configs.CcSecNotifiedServer{}, &common_configs.CcSecNotifiedServerList{}},
		[]string{"ID", "Address", "TsigID"},
		[]string{"{{ .ID }}", "{{ .Address }}", "{{ .TsigID }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&common_configs.CcSecTransferAcl{}, &common_configs.CcSecTransferAclList{}},
		[]string{"ID", "Network", "TsigID"},
		[]string{"{{ .ID }}", "{{ .Network }}", "{{ .TsigID }}"})
}
