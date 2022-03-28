package printer

import (
	. "github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

func init() {
	SetBaseHumanReadablePrinter([]api.Spec{&utils.CommandResults{}},
		[]string{"RequestID", "Status", "ErrorType", "ErrorMessage"},
		[]string{"{{ .RequestID }}", "{{ if .Job }}{{ .Job.Status }}{{ end }}", "{{ if .Job }}{{ .Job.ErrorType }}{{ end }}", "{{ if .Job }}{{ .Job.ErrorMessage }}{{ end }}"})
}
