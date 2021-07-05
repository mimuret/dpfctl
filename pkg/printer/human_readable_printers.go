package printer

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/mimuret/dpfctl/pkg/utils"
)

var (
	humanReadablePrinters = humanReadablePrintersMap{}
	baseArgsTemplateRgexp = regexp.MustCompile("{{([^{}]*)}}")
)

type HumanReadablePrinter interface {
	GetHeaders() []interface{}
	GetRow(interface{}) []interface{}
}

func GetHumanReadablePrinter(name string) HumanReadablePrinter {
	return humanReadablePrinters[name]
}
func SetHumanReadablePrinter(name string, hPrinter HumanReadablePrinter) {
	humanReadablePrinters[name] = hPrinter
}

type humanReadablePrintersMap map[string]HumanReadablePrinter

func SetBaseHumanReadablePrinter(name string, headers []string, argsTemplate []string) {
	//	baseArgsTemplateRgexp := regexp.MustCompile("{{([^{}]*)}}")
	//	matches := baseArgsTemplateRgexp.FindAllStringSubmatch(argsTemplateStr, -1)
	rowTemplate := []*template.Template{}
	for _, argTemplate := range argsTemplate {
		t := template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(argTemplate))
		rowTemplate = append(rowTemplate, t)
	}
	SetHumanReadablePrinter(name, &BaseHumanReadablePrinter{
		Headers:     headers,
		RowTemplate: rowTemplate,
	})
}

type BaseHumanReadablePrinter struct {
	Headers     []string
	RowTemplate []*template.Template
}

func (h *BaseHumanReadablePrinter) GetHeaders() []interface{} {
	return utils.StringSliceToInterfaceSlice(h.Headers)
}

func (h *BaseHumanReadablePrinter) GetRow(i interface{}) []interface{} {
	res := []interface{}{}
	buf := bytes.NewBuffer(nil)
	for _, t := range h.RowTemplate {
		if err := t.Execute(buf, i); err != nil {
			res = append(res, err)
		} else {
			res = append(res, buf.String())
		}
		buf.Reset()
	}
	return res
}

func NewGoTemplateBasePrinter(templateStr string) (*BaseHumanReadablePrinter, error) {
	t, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	return &BaseHumanReadablePrinter{
		Headers:     []string{},
		RowTemplate: []*template.Template{t},
	}, nil
}
