package output

import (
	"fmt"

	"github.com/BaseMax/go-api-contract/pkg/reporter"
	"github.com/fatih/color"
)

var (
	// Color functions
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
	
	// Styled output
	greenBold  = color.New(color.FgGreen, color.Bold).SprintFunc()
	redBold    = color.New(color.FgRed, color.Bold).SprintFunc()
)

// PrintReport prints a formatted report to stdout
func PrintReport(report *reporter.Report) {
	fmt.Println()
	fmt.Println(bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println(bold("  API Contract Validation Report"))
	fmt.Println(bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()

	// Status
	if report.Valid {
		fmt.Printf("%s %s\n", greenBold("✓"), green("All validations passed"))
	} else {
		fmt.Printf("%s %s\n", redBold("✗"), red("Validation failed"))
	}

	// Details
	fmt.Println()
	fmt.Printf("%s: %s\n", cyan("Timestamp"), report.Timestamp)
	fmt.Printf("%s: %d\n", cyan("Status Code"), report.StatusCode)
	fmt.Printf("%s: %d\n", cyan("Errors Found"), report.ErrorCount)

	// Errors
	if len(report.Errors) > 0 {
		fmt.Println()
		fmt.Println(bold("Validation Errors:"))
		fmt.Println()
		
		for i, err := range report.Errors {
			fmt.Printf("  %s. %s\n", yellow(fmt.Sprintf("%d", i+1)), bold(err.Field))
			fmt.Printf("     %s: %s\n", red("Message"), err.Message)
			fmt.Printf("     %s: %s\n", blue("Type"), err.Type)
			fmt.Println()
		}
	}

	// Response body preview
	if report.ResponseBody != nil {
		fmt.Println(bold("Response Body:"))
		fmt.Printf("  %v\n", report.ResponseBody)
		fmt.Println()
	}

	fmt.Println(bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Printf("%s %s\n", greenBold("✓"), green(message))
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Fprintf(color.Output, "%s %s\n", redBold("✗"), red(message))
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Printf("%s %s\n", yellowBold("⚠"), yellow(message))
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Printf("%s %s\n", blueBold("ℹ"), blue(message))
}

var (
	yellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	blueBold   = color.New(color.FgBlue, color.Bold).SprintFunc()
)
