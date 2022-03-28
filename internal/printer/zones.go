package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/zones"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&zones.CurrentRecordList{}},
		[]string{"ID", "Name", "TTL", "RRtype", "RData"},
		[]string{"{{ .ID }}", "{{ .Name }}", "{{ .TTL }}", "{{ .RRType }}", `{{ .RData }}`})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.Dnssec{}},
		[]string{"Enabled", "State", "DsState"},
		[]string{"{{ .Enabled }}", "{{ .State }}", "{{ .DsState }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.DsRecordList{}},
		[]string{"TransitAt", "RDATA"},
		[]string{"{{ .TransitAt }}", "{{ .RRSet }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.ManagedDnsList{}},
		[]string{"ServerName"},
		[]string{"{{ . }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.Record{}, &zones.RecordList{}},
		[]string{"ID", "Name", "TTL", "RRtype", "RData", "State", "Operator"},
		[]string{"{{ .ID }}", "{{ .Name }}", "{{ .TTL }}", "{{ .RRType }}", `{{ .RData  }}`, "{{ .State }}", "{{ .Operator }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.DefaultTTL{}},
		[]string{"Value", "State", "Operator"},
		[]string{"{{ .Value }}", "{{ .State }}", "{{ .Operator }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.HistoryList{}},
		[]string{"ID", "CommittedAt", "Operator", "Description"},
		[]string{"{{ .ID }}", "{{ .CommittedAt }}", "{{ .Operator }}", "{{ .Description }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.LogList{}},
		[]string{"RequestID", "Time", "Status"},
		[]string{"{{ .RequestID }}", "{{ .Time }}", "{{ .Status }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.ZoneProxy{}},
		[]string{"Enabled"},
		[]string{"{{ .Enabled }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.ZoneProxyHealthCheckList{}},
		[]string{"Address", "Status", "TsigName", "Enabled"},
		[]string{"{{ .Address }}", "{{ .Status }}", "{{ .TsigName }}", "{{ .Enabled }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.Contract{}},
		[]string{"ContractID", "ServiceCode", "State"},
		[]string{"{{ .ID }}", "{{ .ServiceCode }}", "{{ .State }}"})

	// zone_default_ttl_diffs
	// record_diffs

}
