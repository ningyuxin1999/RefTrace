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
	results, err := corelint.NFCoreLint(dir)
	if err != nil {
		color.New(color.FgRed).Printf("Error: %s\n", err)
		os.Exit(1)
	}

	errorPrinter := color.New(color.FgRed)
	warningPrinter := color.New(color.FgYellow)
	pathPrinter := color.New(color.FgCyan)

	hasErrors := false

	// Sort results by module path
	sort.Slice(results, func(i, j int) bool {
		return results[i].ModulePath < results[j].ModulePath
	})

	// Print warnings and errors for each module
	for _, result := range results {
		if len(result.Warnings) > 0 || len(result.Errors) > 0 {
			fmt.Printf("\nModule: %s\n", result.ModulePath)
		}

		for _, warning := range result.Warnings {
			pathPrinter.Printf("  • Line %d\n", warning.Line)
			warningPrinter.Printf("    %s\n", warning.Warning)
		}

		for _, err := range result.Errors {
			hasErrors = true
			pathPrinter.Printf("  • Line %d\n", err.Line)
			errorPrinter.Printf("    %s\n", err.Error.Error())
		}
	}

	if hasErrors {
		os.Exit(1)
	}
}
