package reporter

import (
	"testing"

	"github.com/BaseMax/go-api-contract/pkg/validator"
)

func TestGenerateReport(t *testing.T) {
	tests := []struct {
		name   string
		result *validator.ValidationResult
		valid  bool
	}{
		{
			name: "Valid result",
			result: &validator.ValidationResult{
				Valid:      true,
				StatusCode: 200,
				Errors:     []validator.ValidationError{},
			},
			valid: true,
		},
		{
			name: "Invalid result with errors",
			result: &validator.ValidationResult{
				Valid:      false,
				StatusCode: 400,
				Errors: []validator.ValidationError{
					{
						Field:   "field1",
						Message: "error message",
						Type:    "type1",
					},
				},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := GenerateReport(tt.result)

			if report.Valid != tt.valid {
				t.Errorf("GenerateReport().Valid = %v, want %v", report.Valid, tt.valid)
			}

			if report.StatusCode != tt.result.StatusCode {
				t.Errorf("GenerateReport().StatusCode = %v, want %v", report.StatusCode, tt.result.StatusCode)
			}

			if len(report.Errors) != len(tt.result.Errors) {
				t.Errorf("GenerateReport().Errors length = %v, want %v", len(report.Errors), len(tt.result.Errors))
			}

			if report.Timestamp == "" {
				t.Error("GenerateReport().Timestamp is empty")
			}
		})
	}
}

func TestGetSummary(t *testing.T) {
	tests := []struct {
		name     string
		report   *Report
		expected string
	}{
		{
			name: "Valid report",
			report: &Report{
				Valid: true,
			},
			expected: "✓ All validations passed",
		},
		{
			name: "Invalid report",
			report: &Report{
				Valid: false,
			},
			expected: "✗ Validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.report.GetSummary()
			if result != tt.expected {
				t.Errorf("GetSummary() = %v, want %v", result, tt.expected)
			}
		})
	}
}
