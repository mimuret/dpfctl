package printer

import (
	"io"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

type JsonPrinter struct{}

func (p *JsonPrinter) Print(w io.Writer, obj api.Spec) error {
	bs, err := api.MarshalOutput(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}
