package validate

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   Metadata               `yaml:"metadata"`
	Spec       map[string]interface{} `yaml:"spec"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

func ParseManifest(data []byte) (*Manifest, error) {
	var doc Manifest
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func ParseManifests(data []byte) ([]*Manifest, error) {
	var manifests []*Manifest
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var doc Manifest
		err := decoder.Decode(&doc)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		manifests = append(manifests, &doc)
	}
	return manifests, nil
}

func ParseManifestNode(data []byte) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}