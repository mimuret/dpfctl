package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/contracts"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.CommonConfig{}, &contracts.CommonConfigList{}},
		[]string{"ID", "Name", "ManagedDNSEnabled", "Default"},
		[]string{"{{ .ID }}", "{{ .Name }}", "{{ .ManagedDNSEnabled }}", "{{ .Default }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.LogList{}},
		[]string{"RequestID", "Time", "Status"},
		[]string{"{{ .RequestID }}", "{{ .Time }}", "{{ .Status }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.ContractPartnerList{}},
		[]string{"ServiceCode"},
		[]string{"{{ .ServiceCode }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.QpsHistoryList{}},
		[]string{"ServiceCode", "Name", "LastMonth", "LastQps"},
		[]string{"{{ .ServiceCode }}", "{{ .Name }}", "{{ $last := .Values | last}}{{ $last.Month }}", "{{ $last := .Values | last}}{{ $last.Qps }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.ContractZoneList{}},
		[]string{"ZoneID", "ServiceCode", "Name", "State", "CommonConfigID"},
		[]string{"{{ .ID }}", "{{ .ServiceCode }}", "{{ .Name }}", "{{ .State }}", "{{ .CommonConfigID }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.TsigList{}, &contracts.Tsig{}},
		[]string{"ID", "Name", "Algorithm"},
		[]string{"{{ .CommonConfigID }}", "{{ .ID }}", "{{ .Name }}", "{{ .Algorithm }}"})
}
