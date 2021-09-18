package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/zones"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&zones.CurrentRecordList{}},
		[]string{"Id", "Name", "TTL", "RRtype", "RData"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .TTL }}", "{{ .RRType }}", `{{ .RData }}`})
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
		[]string{"Id", "Name", "TTL", "RRtype", "RData", "State", "Operator"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .TTL }}", "{{ .RRType }}", `{{ .RData  }}`, "{{ .State }}", "{{ .Operator }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.DefaultTTL{}},
		[]string{"Value", "State", "Operator"},
		[]string{"{{ .Value }}", "{{ .State }}", "{{ .Operator }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.HistoryList{}},
		[]string{"Id", "CommittedAt", "Operator", "Description"},
		[]string{"{{ .Id }}", "{{ .CommittedAt }}", "{{ .Operator }}", "{{ .Description }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.LogList{}},
		[]string{"RequestId", "Time", "Status"},
		[]string{"{{ .RequestId }}", "{{ .Time }}", "{{ .Status }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.ZoneProxy{}},
		[]string{"Enabled"},
		[]string{"{{ .Enabled }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.ZoneProxyHealthCheckList{}},
		[]string{"Address", "Status", "TsigName", "Enabled"},
		[]string{"{{ .Address }}", "{{ .Status }}", "{{ .TsigName }}", "{{ .Enabled }}"})
	SetBaseHumanReadablePrinter([]api.Spec{&zones.Contract{}},
		[]string{"ContractId", "ServiceCode", "State"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .State }}"})

	// zone_default_ttl_diffs
	// record_diffs

}
