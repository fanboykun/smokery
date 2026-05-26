// Package loader reads project configs, smoke plans, and OpenAPI specs from
// the local filesystem. It supports YAML (preferred for hand-authored config)
// and JSON.
package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// ProjectConfig loads a model.ProjectConfig from a YAML or JSON file.
func ProjectConfig(path string) (*model.ProjectConfig, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	var cfg model.ProjectConfig
	if err := unmarshal(path, data, &cfg); err != nil {
		return nil, fmt.Errorf("parse project config %s: %w", path, err)
	}
	return &cfg, nil
}

// SmokePlan loads a model.SmokePlan from a YAML or JSON file.
func SmokePlan(path string) (*model.SmokePlan, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	var plan model.SmokePlan
	if err := unmarshal(path, data, &plan); err != nil {
		return nil, fmt.Errorf("parse smoke plan %s: %w", path, err)
	}
	return &plan, nil
}

// OpenAPISpec reads an OpenAPI document file (raw bytes for spec.Parse).
func OpenAPISpec(path string) ([]byte, error) {
	return readFile(path)
}

func readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return data, nil
}

func unmarshal(path string, data []byte, out any) error {
	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, out)
	case ".json":
		return json.Unmarshal(data, out)
	default:
		// Try YAML first, then JSON.
		if err := yaml.Unmarshal(data, out); err == nil {
			return nil
		}
		return json.Unmarshal(data, out)
	}
}
