package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"gopkg.in/yaml.v3"
)

func YamlToJson(in []byte) ([]byte, error) {
	mapS := map[string]interface{}{}
	if err := yaml.Unmarshal(in, &mapS); err != nil {
		return nil, err
	}
	out, err := json.Marshal(mapS)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func JsonToYaml(bs json.RawMessage) ([]byte, error) {
	mapS := map[string]interface{}{}
	err := yaml.Unmarshal(bs, &mapS)
	if err != nil {
		return nil, err
	}
	out, err := yaml.Marshal(mapS)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ReadFile(filename string) ([]byte, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".yaml", ".yml":
		bs, err = YamlToJson(bs)
		if err != nil {
			return nil, err
		}
	}
	return bs, nil
}

func ReadSpec(spec apis.Spec, filename string) error {
	bs, err := ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file `%s`: %w", filename, err)
	}
	return api.UnMarshalInput(bs, spec)
}
