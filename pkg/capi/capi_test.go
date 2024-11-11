package main

import (
	"os"
	"testing"
)

func TestProcessAndContainer(t *testing.T) {
	// Create a temporary file with the process content
	content := `
process FOO_BAR {
    tag "$meta.id"
    label 'process_single'

    conda "conda-forge::sed=4.7"
    container "${ workflow.containerEngine == 'singularity' && !task.ext.singularity_pull_docker_container ?
        'https://depot.galaxyproject.org/singularity/ubuntu:20.04' :
        'nf-core/ubuntu:20.04' }"

    input:
    tuple val(meta), path(reads, stageAs: "input*/*")

    output:
    tuple val(meta), path("*.merged.fastq.gz"), emit: reads
    path "versions.yml"                       , emit: versions

    when:
    task.ext.when == null || task.ext.when

    script:
    def args = task.ext.args ?: ''
}`

	// Write to temporary file
	tmpfile := t.TempDir() + "/test.nf"
	if err := os.WriteFile(tmpfile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create module
	moduleHandle, err := createModuleFromFile(tmpfile)
	if err != nil {
		t.Fatalf("Failed to create module: %v", err)
	}
	defer Module_Free(moduleHandle)

	// Get process
	procHandle := Module_GetProcess(moduleHandle, 0)
	if procHandle == 0 {
		t.Fatal("Failed to get process")
	}

	// Check process name
	name := getProcessName(procHandle)
	if name != "FOO_BAR" {
		t.Errorf("Expected process name 'FOO_BAR', got %s", name)
	}

	// Get container directive
	dirCount := Process_GetDirectiveCount(procHandle)
	var containerFound bool
	for i := 0; i < int(dirCount); i++ {
		dirHandle := Process_GetDirective(procHandle, cInt(i))
		if dirHandle == 0 {
			continue
		}
		if Directive_IsContainer(dirHandle) == 1 {
			containerFound = true
			format := getContainerFormat(dirHandle)
			if format != "ternary" {
				t.Errorf("Expected container format 'ternary', got %s", format)
			}

			condition := getContainerCondition(dirHandle)
			expectedCond := "((workflow.containerEngine == singularity) && !(task.ext.singularity_pull_docker_container))"
			if condition != expectedCond {
				t.Errorf("Expected condition '%s', got '%s'", expectedCond, condition)
			}

			trueName := getContainerTrueName(dirHandle)
			expectedTrue := "https://depot.galaxyproject.org/singularity/ubuntu:20.04"
			if trueName != expectedTrue {
				t.Errorf("Expected true name '%s', got '%s'", expectedTrue, trueName)
			}

			falseName := getContainerFalseName(dirHandle)
			expectedFalse := "nf-core/ubuntu:20.04"
			if falseName != expectedFalse {
				t.Errorf("Expected false name '%s', got '%s'", expectedFalse, falseName)
			}
		}
	}

	if !containerFound {
		t.Error("No container directive found")
	}
}
