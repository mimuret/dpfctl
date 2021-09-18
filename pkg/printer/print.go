package printer

import (
	"io"
	"strings"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

type Printer interface {
	Print(w io.Writer, obj api.Spec) error
}

func GetPrinter(s api.Spec, outType string) (Printer, error) {
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
	return NewLinePrinter(), nil
}
