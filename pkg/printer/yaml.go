package printer

import (
	"io"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"sigs.k8s.io/yaml"
)

type YAMLPrinter struct{}

func (p *YAMLPrinter) Print(w io.Writer, obj api.Spec) error {
	bs, err := api.MarshalOutput(obj)
	if err != nil {
		return err
	}
	bs, err = yaml.JSONToYAML(bs)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err

}
