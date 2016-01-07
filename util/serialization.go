package util

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Format is a (de)serialization format
type Format int

const (
	// JSON format
	JSON Format = iota
	// YAML format
	YAML
)

// Marshal serializes o into a string using the format f.
func Marshal(o interface{}, f Format) (string, error) {
	switch f {
	case YAML:
		bytes, err := yaml.Marshal(o)
		if err != nil {
			return "", err
		}
		return string(bytes), nil

	case JSON:
		bytes, err := json.Marshal(o)
		if err != nil {
			return "", err
		}
		return string(bytes), nil

	default:
		return "", fmt.Errorf("unsupported format %v", f)
	}
}

// Unmarshal deserializes i into o using the format f.
func Unmarshal(i []byte, o interface{}, f Format) error {
	switch f {
	case YAML:
		return yaml.Unmarshal(i, o)

	case JSON:
		return json.Unmarshal(i, o)

	default:
		return fmt.Errorf("unsupported format %v", f)
	}
}
