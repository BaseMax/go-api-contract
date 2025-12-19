package loader

import (
	"context"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

// LoadSpec loads an OpenAPI specification from a file
func LoadSpec(filePath string) (*openapi3.T, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("spec file not found: %s", filePath)
	}

	// Load the spec
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	spec, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec: %w", err)
	}

	// Validate the spec
	ctx := context.Background()
	if err := spec.Validate(ctx); err != nil {
		return nil, fmt.Errorf("invalid spec: %w", err)
	}

	return spec, nil
}
