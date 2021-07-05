package printer

func init() {
	SetBaseHumanReadablePrinter("contracts",
		[]string{"ContractId", "ServiceCode", "State"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .State }}"})
	SetBaseHumanReadablePrinter("delegations",
		[]string{"ZoneId", "ServiceCode", "Name", "Network"},
		[]string{"{{ .Id }}", "{{ .ServiceCode }}", "{{ .Name }}", "{{ .Network }}"})
	SetBaseHumanReadablePrinter("jobs",
		[]string{"JobId", "Status", "ErrorType", "ResourceUrl", "ErrorType", "ErrorMessage"},
		[]string{"{{ .RequestId }}", "{{ .Status }}", "{{ .ResourceUrl }}", "{{ .ErrorType }}", "{{ .ErrorMessage }}"})
	SetBaseHumanReadablePrinter("zones",
		[]string{"ZoneId", "Name", "State", "CommonConfigId"},
		[]string{"{{ .Id }}", "{{ .Name }}", "{{ .State }}", "{{ .CommonConfigId }}"})
}
