package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/common_configs"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&common_configs.CcPrimary{}, &common_configs.CcPrimaryList{}},
		[]string{"Id", "Address", "TsigId", "Enabled"},
		[]string{"{{ .Id }}", "{{ .Address }}", "{{ .TsigId }}", "{{ .Enabled }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&common_configs.CcSecNotifiedServer{}, &common_configs.CcSecNotifiedServerList{}},
		[]string{"Id", "Address", "TsigId"},
		[]string{"{{ .Id }}", "{{ .Address }}", "{{ .TsigId }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&common_configs.CcSecTransferAcl{}, &common_configs.CcSecTransferAclList{}},
		[]string{"Id", "Network", "TsigId"},
		[]string{"{{ .Id }}", "{{ .Network }}", "{{ .TsigId }}"})
}
