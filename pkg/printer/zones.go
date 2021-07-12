package printer

func init() {
	SetBaseHumanReadablePrinter("current_records",
		[]string{"Id", "Name", "TTL", "RRtype", "RData"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .TTL }}", "{{ .RRType }}", `{{ .RData }}`})
	SetBaseHumanReadablePrinter("dnssec",
		[]string{"Enabled", "State", "DsState"},
		[]string{"{{ .Enabled }}", "{{ .State }}", "{{ .DsState }}"})
	SetBaseHumanReadablePrinter("ds_records",
		[]string{"RRSet", "TransitAt"},
		[]string{"{{ .RRSet }}", "{{ .TransitAt }}"})
	SetBaseHumanReadablePrinter("managed_dns_servers",
		[]string{"Managed DNS Server Name"},
		[]string{"{{ . }}"})
	SetBaseHumanReadablePrinter("records",
		[]string{"Id", "Name", "TTL", "RRtype", "RData", "State", "Operator"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .TTL }}", "{{ .RRType }}", `{{ .RData  }}`, "{{ .State }}", "{{ .Operator }}"})
	SetBaseHumanReadablePrinter("zone_default_ttl",
		[]string{"Value", "State", "Operator"},
		[]string{"{{ .Value }}", "{{ .State }}", "{{ .Operator }}"})
	SetBaseHumanReadablePrinter("zone_histories",
		[]string{"Id", "CommittedAt", "Operator", "Description"},
		[]string{"{{ .Id }}", "{{ .CommittedAt }}", "{{ .Operator }}", "{{ .Description }}"})
	SetBaseHumanReadablePrinter("zone_logs",
		[]string{"RequestId", "Time", "Status"},
		[]string{"{{ .RequestId }}", "{{ .Time }}", "{{ .Status }}"})
	SetBaseHumanReadablePrinter("zone_proxy",
		[]string{"Enabled"},
		[]string{"{{ .Enabled }}"})
	SetBaseHumanReadablePrinter("zone_proxy_health_check",
		[]string{"Address", "Status", "TsigName", "Enabled"},
		[]string{"{{ .Address }}", "{{ .Status }}", "{{ .TsigName }}", "{{ .Enabled }}"})
	SetBaseHumanReadablePrinter("zones_contract",
		[]string{"ContractId", "ServiceCode", "State"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .State }}"})

	// zone_default_ttl_diffs
	// record_diffs

}
