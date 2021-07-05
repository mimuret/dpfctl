package printer

func init() {
	SetBaseHumanReadablePrinter("cc_primaries",
		[]string{"Id", "Address", "TsigId", "Enabled"},
		[]string{"{{ .Id }}", "{{ .Address }}", "{{ .TsigId }}", "{{ .Enabled }}"})
	SetBaseHumanReadablePrinter("cc_sec_notified_servers",
		[]string{"Id", "Address", "TsigId"},
		[]string{"{{ .Id }}", "{{ .Address }}", "{{ .TsigId }}"})
	SetBaseHumanReadablePrinter("cc_sec_transfer_acls",
		[]string{"Id", "Network", "TsigId"},
		[]string{"{{ .Id }}", "{{ .Network }}", "{{ .TsigId }}"})
}
