package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&core.Contract{}, &core.ContractList{}},
		[]string{"ContractID", "ServiceCode", "State"},
		[]string{"{{ .ID }}", "{{ .ServiceCode }}", "{{ .State }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&core.DelegationList{}},
		[]string{"ZoneID", "ServiceCode", "Name", "Network", "LastRequestTime"},
		[]string{"{{ .ID }}", "{{ .ServiceCode }}", "{{ .Name }}", "{{ .Network }}", "{{ .DelegationRequestedAt }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&core.Job{}},
		[]string{"RequestID", "Status", "ErrorType", "ErrorMessage"},
		[]string{"{{ .RequestID }}", "{{ .Status }}", "{{ .ErrorType }}", "{{ .ErrorMessage }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&core.Zone{}, &core.ZoneList{}},
		[]string{"ZoneID", "ServiceCode", "Name", "State", "CommonConfigID"},
		[]string{"{{ .ID }}", "{{ .ServiceCode }}", "{{ .Name }}", "{{ .State }}", "{{ .CommonConfigID }}"})
}
