package validate

import (
	"fmt"
)

func Validate(data []byte) ([]*Manifest, error) {
	manifests, err := ParseManifests(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifests: %w", err)
	}
	return manifests, nil
}