package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/mimuret/golang-iij-dpf/pkg/schema"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

var DefaultFS = afero.NewOsFs()

type ResourceReader struct {
	fs afero.Fs
}

func NewResourceReader(fs afero.Fs) *ResourceReader {
	if fs == nil {
		fs = DefaultFS
	}
	return &ResourceReader{fs: fs}
}

func (reader *ResourceReader) GetResources(filename string) ([]apis.Spec, error) {
	var (
		r   io.Reader
		err error
	)
	if filename == "-" {
		r = os.Stdin
	} else {
		f, err := reader.fs.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	docs, err := reader.ReadYamlDocuments(r)
	if err != nil {
		return nil, err
	}
	return reader.ParseResouress(docs)
}

func (*ResourceReader) ReadYamlDocuments(r io.Reader) ([]json.RawMessage, error) {
	dec := yaml.NewDecoder(r)
	res := []json.RawMessage{}
LOOP:
	for {
		tmp := map[string]interface{}{}
		if err := dec.Decode(&tmp); err != nil {
			if err == io.EOF {
				break LOOP
			}
			return nil, err
		}
		out, err := json.Marshal(tmp)
		if err != nil {
			return nil, err
		}
		res = append(res, out)
	}
	return res, nil
}

func (reader *ResourceReader) ParseResouress(raws []json.RawMessage) ([]apis.Spec, error) {
	res := []apis.Spec{}
	for i, raw := range raws {
		s, err := schema.SchemaSet.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("document[%d] is failed to parse: %w", i, err)
		}
		res = append(res, s)
	}
	return res, nil
}
