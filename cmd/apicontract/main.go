package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BaseMax/go-api-contract/pkg/client"
	"github.com/BaseMax/go-api-contract/pkg/loader"
	"github.com/BaseMax/go-api-contract/pkg/output"
	"github.com/BaseMax/go-api-contract/pkg/reporter"
	"github.com/BaseMax/go-api-contract/pkg/validator"
)

const version = "1.0.0"

func main() {
	// Define flags
	specFile := flag.String("spec", "", "Path to OpenAPI/Swagger specification file (required)")
	endpoint := flag.String("endpoint", "", "API endpoint to validate (required)")
	method := flag.String("method", "GET", "HTTP method to use")
	body := flag.String("body", "", "Request body (for POST/PUT/PATCH)")
	headers := flag.String("headers", "", "Request headers as JSON object")
	jsonOutput := flag.Bool("json", false, "Output results as JSON instead of colored text")
	verboseMode := flag.Bool("verbose", false, "Enable verbose output")
	versionFlag := flag.Bool("version", false, "Print version information")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("apicontract version %s\n", version)
		os.Exit(0)
	}

	if *specFile == "" || *endpoint == "" {
		fmt.Fprintf(os.Stderr, "Error: -spec and -endpoint are required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Load OpenAPI spec
	spec, err := loader.LoadSpec(*specFile)
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to load spec: %v", err))
		os.Exit(1)
	}

	if *verboseMode {
		output.PrintInfo(fmt.Sprintf("Loaded spec: %s", *specFile))
	}

	// Parse headers
	var headerMap map[string]string
	if *headers != "" {
		if err := json.Unmarshal([]byte(*headers), &headerMap); err != nil {
			output.PrintError(fmt.Sprintf("Failed to parse headers: %v", err))
			os.Exit(1)
		}
	}

	// Make HTTP request
	httpClient := client.NewClient()
	response, err := httpClient.Do(*method, *endpoint, *body, headerMap)
	if err != nil {
		output.PrintError(fmt.Sprintf("Failed to make request: %v", err))
		os.Exit(1)
	}

	if *verboseMode {
		output.PrintInfo(fmt.Sprintf("Request: %s %s", *method, response.URL))
		output.PrintInfo(fmt.Sprintf("Response Status: %d", response.StatusCode))
	}

	// Extract path from actual URL used in the request
	endpointPath := response.URL
	if idx := strings.Index(endpointPath, "://"); idx != -1 {
		// Remove protocol
		endpointPath = endpointPath[idx+3:]
		// Remove host
		if idx := strings.Index(endpointPath, "/"); idx != -1 {
			endpointPath = endpointPath[idx:]
		} else {
			endpointPath = "/"
		}
	}

	// Validate response
	v := validator.NewValidator(spec)
	result, err := v.Validate(*method, endpointPath, response)
	if err != nil {
		output.PrintError(fmt.Sprintf("Validation error: %v", err))
		os.Exit(1)
	}

	// Generate report
	report := reporter.GenerateReport(result)

	// Output results
	if *jsonOutput {
		jsonData, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			output.PrintError(fmt.Sprintf("Failed to marshal JSON: %v", err))
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	} else {
		output.PrintReport(report)
	}

	// Exit with appropriate code
	if !report.Valid {
		os.Exit(1)
	}
}
