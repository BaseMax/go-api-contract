package reporter

import (
	"time"

	"github.com/BaseMax/go-api-contract/pkg/validator"
)

// Report represents the validation report
type Report struct {
	Valid        bool                        `json:"valid"`
	StatusCode   int                         `json:"statusCode"`
	Errors       []validator.ValidationError `json:"errors,omitempty"`
	ResponseBody interface{}                 `json:"responseBody,omitempty"`
	Timestamp    string                      `json:"timestamp"`
	ErrorCount   int                         `json:"errorCount"`
}

// GenerateReport generates a report from validation results
func GenerateReport(result *validator.ValidationResult) *Report {
	return &Report{
		Valid:        result.Valid,
		StatusCode:   result.StatusCode,
		Errors:       result.Errors,
		ResponseBody: result.ResponseBody,
		Timestamp:    time.Now().Format(time.RFC3339),
		ErrorCount:   len(result.Errors),
	}
}

// GetSummary returns a summary string of the report
func (r *Report) GetSummary() string {
	if r.Valid {
		return "✓ All validations passed"
	}
	return "✗ Validation failed"
}
