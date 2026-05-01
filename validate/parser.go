package validate

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

func SplitDocuments(data []byte) ([][]byte, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var docs [][]byte
	for {
		var node yaml.Node
		err := decoder.Decode(&node)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if len(node.Content) == 0 {
			continue
		}
		buf, err := yaml.Marshal(node.Content[0])
		if err != nil {
			return nil, err
		}
		docs = append(docs, buf)
	}
	return docs, nil
}
