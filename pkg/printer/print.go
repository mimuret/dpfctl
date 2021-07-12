package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

type Printer interface {
	Print(w io.Writer, obj api.Spec) error
}

func GetPrinter(cmdName, outType string) (Printer, error) {
	var arg string
	outputType := strings.SplitN(outType, "=", 2)
	if len(outputType) == 2 {
		arg = outputType[1]
	}
	switch outputType[0] {
	case "yaml":
		return &YAMLPrinter{}, nil
	case "json":
		return &JsonPrinter{}, nil
	case "go-template":
		return NewGoTemplatePrinter(arg)
	}
	return NewLinePrinter(cmdName)
}

type YAMLPrinter struct{}

func (p *YAMLPrinter) Print(w io.Writer, obj api.Spec) error {
	bs, err := api.MarshalOutput(obj)
	if err != nil {
		return err
	}
	bs, err = utils.JsonToYaml(bs)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err

}

type JsonPrinter struct{}

func (p *JsonPrinter) Print(w io.Writer, obj api.Spec) error {
	bs, err := api.MarshalOutput(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}

func NewGoTemplatePrinter(templateStr string) (*LinePrinter, error) {
	printer, err := NewGoTemplateBasePrinter(templateStr)
	if err != nil {
		return nil, err
	}
	return &LinePrinter{Printer: printer}, nil
}

type LinePrinter struct {
	Noheaders bool
	Printer   HumanReadablePrinter
}

func (p *LinePrinter) SetNoHeaders(s bool) {
	p.Noheaders = s
}

func NewLinePrinter(cmdName string) (*LinePrinter, error) {
	printer := GetHumanReadablePrinter(cmdName)

	if printer == nil {
		return nil, fmt.Errorf("not support print %s", cmdName)
	}
	return &LinePrinter{Printer: printer}, nil
}

func (p *LinePrinter) Print(w io.Writer, obj api.Spec) error {
	table := uitable.New()

	if !p.Noheaders {
		table.AddRow(p.Printer.GetHeaders()...)
	}
	if l, ok := obj.(api.ListSpec); ok {
		for i := 0; i < l.Len(); i++ {
			table.AddRow(p.Printer.GetRow(l.Index(i))...)
		}
	} else {
		table.AddRow(p.Printer.GetRow(obj)...)
	}
	_, err := w.Write(table.Bytes())
	w.Write([]byte("\n"))
	return err
}
