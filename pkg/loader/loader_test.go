package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSpec(t *testing.T) {
	// Create a temporary valid OpenAPI spec
	validSpec := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: Success
`

	tmpDir := t.TempDir()
	validSpecFile := filepath.Join(tmpDir, "valid.yaml")

	if err := os.WriteFile(validSpecFile, []byte(validSpec), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Valid spec",
			file:    validSpecFile,
			wantErr: false,
		},
		{
			name:    "File not found",
			file:    filepath.Join(tmpDir, "notfound.yaml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := LoadSpec(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && spec == nil {
				t.Error("LoadSpec() returned nil spec for valid file")
			}
		})
	}
}
