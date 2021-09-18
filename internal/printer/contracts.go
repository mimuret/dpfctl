package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/contracts"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.CommonConfig{}, &contracts.CommonConfigList{}},
		[]string{"Id", "Name", "ManagedDNSEnabled", "Default"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .ManagedDNSEnabled }}", "{{ .Default }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.LogList{}},
		[]string{"RequestId", "Time", "Status"},
		[]string{"{{ .RequestId }}", "{{ .Time }}", "{{ .Status }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.ContractPartnerList{}},
		[]string{"ServiceCode"},
		[]string{"{{ .ServiceCode }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.QpsHistoryList{}},
		[]string{"ServiceCode", "Name", "LastMonth", "LastQps"},
		[]string{"{{ .ServiceCode }}", "{{ .Name }}", "{{ $last := .Values | last}}{{ $last.Month }}", "{{ $last := .Values | last}}{{ $last.Qps }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.ContractZoneList{}},
		[]string{"ZoneId", "ServiceCode", "Name", "State", "CommonConfigId"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .Name }}", "{{ .State }}", "{{ .CommonConfigId }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&contracts.TsigList{}, &contracts.Tsig{}},
		[]string{"Id", "Name", "Algorithm"},
		[]string{"{{ .CommonConfigId }}", "{{ .Id }}", "{{ .Name }}", "{{ .Algorithm }}"})
}
