package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&core.Contract{}, &core.ContractList{}},
		[]string{"ContractId", "ServiceCode", "State"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .State }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&core.DelegationList{}},
		[]string{"ZoneId", "ServiceCode", "Name", "Network", "LastRequestTime"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .Name }}", "{{ .Network }}", "{{ .DelegationRequestedAt }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&core.Job{}},
		[]string{"RequestId", "Status", "ErrorType", "ErrorMessage"},
		[]string{"{{ .RequestId }}", "{{ .Status }}", "{{ .ErrorType }}", "{{ .ErrorMessage }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&core.Zone{}, &core.ZoneList{}},
		[]string{"ZoneId", "ServiceCode", "Name", "State", "CommonConfigId"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .Name }}", "{{ .State }}", "{{ .CommonConfigId }}"})
}
