package printer

func init() {
	SetBaseHumanReadablePrinter("common_configs",
		[]string{"Id", "Name", "ManagedDNSEnabled", "Default"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .ManagedDNSEnabled }}", "{{ .Default }}"})
	SetBaseHumanReadablePrinter("contract_logs",
		[]string{"RequestId", "Time", "Status"},
		[]string{"{{ .RequestId }}", "{{ .Time }}", "{{ .Status }}"})
	SetBaseHumanReadablePrinter("contract_partners",
		[]string{"ServiceCode"},
		[]string{"{{ .ServiceCode }}"})
	SetBaseHumanReadablePrinter("contract_qps",
		[]string{"ServiceCode", "Name", "LastMonth", "LastQps"},
		[]string{"{{ .ServiceCode }}", "{{ .Name }}", "{{ last .Values | .Month	 }}", "{{ last .Values | .Qps }}"})
	SetBaseHumanReadablePrinter("contract_zone",
		[]string{"ZoneId", "Name", "State", "CommonConfigId"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .State }}", "{{ .CommonConfigId }}"})
	SetBaseHumanReadablePrinter("tsigs",
		[]string{"Id", "Name", "Algorithm"},
		[]string{"{{ .CommonConfigId }}", "{{ .Id }}", "{{ .Name }}", "{{ .Algorithm }}"})
	SetBaseHumanReadablePrinter("tsigs_common_configs",
		[]string{"Id", "Name", "ManagedDNSEnabled", "Default"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .ManagedDNSEnabled }}", "{{ .Default }}"})
}
