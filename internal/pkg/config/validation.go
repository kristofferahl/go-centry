package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
)

func validateManifestYaml(schema string, r io.Reader) error {
	jsonschema.Loaders["bindata"] = loadFileBinData

	s, err := jsonschema.Compile(schema)
	if err != nil {
		return err
	}

	if err = s.Validate(r); err != nil {
		m := strings.Replace(err.Error(), "bindata://schemas/manifest.json#", "schemas/manifest.json", 1)
		return fmt.Errorf("invalid manifest file\n\n%s", m)
	}

	return nil
}

func loadFileBinData(s string) (io.ReadCloser, error) {
	s = strings.Replace(s, "bindata://", "", -1)

	data, err := Asset(s)
	if err != nil {
		return nil, fmt.Errorf("failed get the manifest schema, %s", err.Error())
	}

	return ioutil.NopCloser(bytes.NewReader(data)), nil
}