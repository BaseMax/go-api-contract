package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/BaseMax/go-api-contract/pkg/client"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// ValidationResult holds the result of a validation
type ValidationResult struct {
	Valid        bool
	StatusCode   int
	Errors       []ValidationError
	ResponseBody interface{}
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
	Type    string
}

// Validator validates API responses against OpenAPI specs
type Validator struct {
	spec   *openapi3.T
	router routers.Router
}

// NewValidator creates a new validator
func NewValidator(spec *openapi3.T) *Validator {
	router, err := gorillamux.NewRouter(spec)
	if err != nil {
		// Fallback: if gorilla router fails, we'll handle it in validation
		router = nil
	}

	return &Validator{
		spec:   spec,
		router: router,
	}
}

// Validate validates an API response against the spec
func (v *Validator) Validate(method, path string, response *client.Response) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:      true,
		StatusCode: response.StatusCode,
		Errors:     []ValidationError{},
	}

	// Parse response body
	if len(response.Body) > 0 {
		var body interface{}
		if err := json.Unmarshal(response.Body, &body); err == nil {
			result.ResponseBody = body
		} else {
			result.ResponseBody = string(response.Body)
		}
	}

	// Find the operation in the spec
	_, operation, err := v.findOperation(method, path)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "path",
			Message: err.Error(),
			Type:    "path_not_found",
		})
		return result, nil
	}

	// Validate status code
	if !v.validateStatusCode(operation, response.StatusCode) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "statusCode",
			Message: fmt.Sprintf("status code %d is not defined in the spec", response.StatusCode),
			Type:    "invalid_status_code",
		})
	}

	// Validate response schema
	if err := v.validateResponseSchema(operation, response); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "response",
			Message: err.Error(),
			Type:    "schema_validation_error",
		})
	}

	// Validate response headers
	if errs := v.validateResponseHeaders(operation, response); len(errs) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, errs...)
	}

	return result, nil
}

// findOperation finds the operation in the spec for the given method and path
func (v *Validator) findOperation(method, urlPath string) (*openapi3.PathItem, *openapi3.Operation, error) {
	// Normalize method
	method = strings.ToLower(method)

	// Try to find exact match first
	for path, pathItem := range v.spec.Paths.Map() {
		if matchPath(path, urlPath) {
			operation := getOperationByMethod(pathItem, method)
			if operation != nil {
				return pathItem, operation, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("operation not found: %s %s", strings.ToUpper(method), urlPath)
}

// matchPath checks if a URL path matches an OpenAPI path pattern
func matchPath(pattern, path string) bool {
	// Simple path matching - convert {param} to * for basic matching
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], "{") && strings.HasSuffix(patternParts[i], "}") {
			continue // Parameter - matches anything
		}
		if patternParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}

// getOperationByMethod gets the operation for a given HTTP method
func getOperationByMethod(pathItem *openapi3.PathItem, method string) *openapi3.Operation {
	switch method {
	case "get":
		return pathItem.Get
	case "post":
		return pathItem.Post
	case "put":
		return pathItem.Put
	case "delete":
		return pathItem.Delete
	case "patch":
		return pathItem.Patch
	case "head":
		return pathItem.Head
	case "options":
		return pathItem.Options
	default:
		return nil
	}
}

// validateStatusCode checks if the status code is defined in the spec
func (v *Validator) validateStatusCode(operation *openapi3.Operation, statusCode int) bool {
	if operation.Responses == nil {
		return true
	}

	// Check for exact match
	if operation.Responses.Status(statusCode) != nil {
		return true
	}

	// Check for wildcard patterns (2XX, 3XX, etc.)
	wildcardPattern := fmt.Sprintf("%dXX", statusCode/100)
	if operation.Responses.Value(wildcardPattern) != nil {
		return true
	}

	// Check for default response
	if operation.Responses.Default() != nil {
		return true
	}

	return false
}

// validateResponseSchema validates the response body against the schema
func (v *Validator) validateResponseSchema(operation *openapi3.Operation, response *client.Response) error {
	if operation.Responses == nil {
		return nil
	}

	responseRef := operation.Responses.Status(response.StatusCode)
	if responseRef == nil {
		// Try wildcard pattern
		wildcardPattern := fmt.Sprintf("%dXX", response.StatusCode/100)
		responseRef = operation.Responses.Value(wildcardPattern)
	}
	if responseRef == nil {
		responseRef = operation.Responses.Default()
	}

	if responseRef == nil || responseRef.Value == nil {
		return nil
	}

	// Get content type
	contentType := "application/json"
	if ct := response.Headers.Get("Content-Type"); ct != "" {
		contentType = strings.Split(ct, ";")[0]
	}

	// Get media type
	content := responseRef.Value.Content
	if content == nil {
		return nil
	}

	mediaType := content.Get(contentType)
	if mediaType == nil {
		return nil
	}

	if mediaType.Schema == nil || mediaType.Schema.Value == nil {
		return nil
	}

	// Parse response body as JSON
	var body interface{}
	if len(response.Body) > 0 {
		if err := json.Unmarshal(response.Body, &body); err != nil {
			return fmt.Errorf("failed to parse response body as JSON: %w", err)
		}
	}

	// Validate schema
	if err := mediaType.Schema.Value.VisitJSON(body); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

// validateResponseHeaders validates response headers
func (v *Validator) validateResponseHeaders(operation *openapi3.Operation, response *client.Response) []ValidationError {
	var errors []ValidationError

	if operation.Responses == nil {
		return errors
	}

	responseRef := operation.Responses.Status(response.StatusCode)
	if responseRef == nil {
		wildcardPattern := fmt.Sprintf("%dXX", response.StatusCode/100)
		responseRef = operation.Responses.Value(wildcardPattern)
	}
	if responseRef == nil {
		responseRef = operation.Responses.Default()
	}

	if responseRef == nil || responseRef.Value == nil {
		return errors
	}

	// Check required headers
	for headerName, headerRef := range responseRef.Value.Headers {
		if headerRef.Value != nil && headerRef.Value.Required {
			if response.Headers.Get(headerName) == "" {
				errors = append(errors, ValidationError{
					Field:   "headers." + headerName,
					Message: fmt.Sprintf("required header '%s' is missing", headerName),
					Type:    "missing_required_header",
				})
			}
		}
	}

	return errors
}

// ValidateWithFilter validates using openapi3filter (alternative method)
func (v *Validator) ValidateWithFilter(method, path string, response *client.Response) error {
	if v.router == nil {
		return fmt.Errorf("router not initialized")
	}

	// Parse the URL path
	parsedURL, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Find route
	route, pathParams, err := v.router.FindRoute(&http.Request{
		Method: method,
		URL:    parsedURL,
	})
	if err != nil {
		return fmt.Errorf("route not found: %w", err)
	}

	// Create response validation input
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    &http.Request{Method: method},
			PathParams: pathParams,
			Route:      route,
		},
		Status: response.StatusCode,
		Header: response.Headers,
	}

	// Set body
	if len(response.Body) > 0 {
		responseValidationInput.SetBodyBytes(response.Body)
	}

	// Validate
	if err := openapi3filter.ValidateResponse(nil, responseValidationInput); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
