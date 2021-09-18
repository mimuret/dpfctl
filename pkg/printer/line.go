package printer

import (
	"io"

	"github.com/gosuri/uitable"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

type LinePrinter struct {
	Noheaders bool
	Printer   HumanReadablePrinter
}

func NewGoTemplatePrinter(templateStr string) (*LinePrinter, error) {
	printer, err := NewGoTemplateBasePrinter(templateStr)
	if err != nil {
		return nil, err
	}
	return &LinePrinter{Printer: printer}, nil
}

func NewLinePrinter() *LinePrinter {
	return &LinePrinter{}
}

func (p *LinePrinter) SetNoHeaders(s bool) {
	p.Noheaders = s
}

func (p *LinePrinter) Print(w io.Writer, obj api.Spec) error {
	table := uitable.New()
	if p.Printer == nil {
		p.Printer = GetHumanReadablePrinter(obj)
	}
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
