package printer

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

var (
	defaultHumanReadablePrinter = DefaultHumanReadablePrinter{}
	humanReadablePrinters       = map[string]humanReadablePrintersMap{}
	baseArgsTemplateRgexp       = regexp.MustCompile("{{([^{}]*)}}")
)

type HumanReadablePrinter interface {
	GetHeaders() []interface{}
	GetRow(interface{}) []interface{}
}

type DefaultHumanReadablePrinter struct{}

func (p *DefaultHumanReadablePrinter) GetHeaders() []interface{} {
	return []interface{}{"not support line print"}
}
func (p *DefaultHumanReadablePrinter) GetRow(v interface{}) []interface{} {
	return []interface{}{v}
}

func GetHumanReadablePrinter(s api.Spec) HumanReadablePrinter {
	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Ptr {
		name := st.Elem().Name()
		panic(fmt.Sprintf("SetHumanReadablePrinter.Add name %s: is not ptr %v", name, s))
	}
	name := st.Elem().Name()
	if _, ok := humanReadablePrinters[s.GetGroup()]; !ok {
		return &DefaultHumanReadablePrinter{}
	}
	if _, ok := humanReadablePrinters[s.GetGroup()][name]; !ok {
		return &DefaultHumanReadablePrinter{}
	}
	return humanReadablePrinters[s.GetGroup()][name]
}

func SetHumanReadablePrinter(s api.Spec, hPrinter HumanReadablePrinter) {
	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Ptr {
		name := st.Elem().Name()
		panic(fmt.Sprintf("SetHumanReadablePrinter.Add name %s: is not ptr %v", name, s))
	}
	name := st.Elem().Name()
	if _, ok := humanReadablePrinters[s.GetGroup()]; !ok {
		humanReadablePrinters[s.GetGroup()] = humanReadablePrintersMap{}
	}
	humanReadablePrinters[s.GetGroup()][name] = hPrinter
}

type humanReadablePrintersMap map[string]HumanReadablePrinter

func SetBaseHumanReadablePrinter(specs []api.Spec, headers []string, argsTemplate []string) {
	//	baseArgsTemplateRgexp := regexp.MustCompile("{{([^{}]*)}}")
	//	matches := baseArgsTemplateRgexp.FindAllStringSubmatch(argsTemplateStr, -1)
	rowTemplate := []*template.Template{}
	for _, argTemplate := range argsTemplate {
		t := template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(argTemplate))
		rowTemplate = append(rowTemplate, t)
	}
	for _, s := range specs {
		SetHumanReadablePrinter(s, &BaseHumanReadablePrinter{
			Headers:     headers,
			RowTemplate: rowTemplate,
		})
	}
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
