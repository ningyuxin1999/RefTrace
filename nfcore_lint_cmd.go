package main

import (
	"fmt"
	"os"
	"reft-go/nf/corelint"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var nfcoreCmd = &cobra.Command{
	Use:   "nfcore",
	Short: "Run nf-core linting rules on a Nextflow pipeline",
	Run:   runNFCoreLint,
}

func init() {
	lintCmd.AddCommand(nfcoreCmd)
	nfcoreCmd.Flags().StringVarP(&dir, "directory", "d", ".", "Directory to lint")
}

func runNFCoreLint(cmd *cobra.Command, args []string) {
	results := corelint.NFCoreLint(dir)

	errorPrinter := color.New(color.FgRed)
	warningPrinter := color.New(color.FgYellow)
	pathPrinter := color.New(color.FgCyan)

	hasErrors := false

	// Sort warnings by module path
	sort.Slice(results.Warnings, func(i, j int) bool {
		return results.Warnings[i].ModulePath < results.Warnings[j].ModulePath
	})

	// Sort errors by module path
	sort.Slice(results.Errors, func(i, j int) bool {
		return results.Errors[i].ModulePath < results.Errors[j].ModulePath
	})

	if len(results.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range results.Warnings {
			pathPrinter.Printf("  • %s:%d\n", warning.ModulePath, warning.Line)
			warningPrinter.Printf("    %s\n", warning.Warning)
		}
	}

	if len(results.Errors) > 0 {
		hasErrors = true
		fmt.Println("\nErrors:")
		for _, err := range results.Errors {
			pathPrinter.Printf("  • %s:%d\n", err.ModulePath, err.Line)
			errorPrinter.Printf("    %s\n", err.Error.Error())
		}
	}

	if hasErrors {
		os.Exit(1)
	}
}
