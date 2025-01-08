from typing import List
import pytest
import tempfile
from pathlib import Path
from reftrace import Module, parse_modules, ModuleResult
from reftrace.directives import ContainerFormat
import os

def test_process_and_container():
    # Create temporary file with the same process content
    content = """
process CAT_FASTQ {
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
}
"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.nf') as tmp:
        tmp.write(content)
        tmp.flush()
        
        # Create module from file
        module_result = Module.from_file(tmp.name)
        assert module_result.error is None
        module = module_result.module
        
        # Should have one process
        assert len(module.processes) == 1
        process = module.processes[0]
        
        # Check process name
        assert process.name == "CAT_FASTQ"

        containers = process.containers
        
        # Check container directive
        containers = [c for c in containers if c.format == ContainerFormat.TERNARY]
        assert len(containers) == 1
        container = containers[0]
        
        # Verify container properties
        assert container.condition == "((workflow.containerEngine == singularity) && !(task.ext.singularity_pull_docker_container))"
        assert container.true_name == "https://depot.galaxyproject.org/singularity/ubuntu:20.04"
        assert container.false_name == "nf-core/ubuntu:20.04"
        
        # Add label directive testing
        labels = process.labels
        assert len(labels) == 1
        label = labels[0]
        assert label.label == "process_single"


def test_workflow_structure():
    content = """
class MyClass {
    // Invalid - missing return type and modifiers
    myMethod() {
        println "Hello"
    }
}
"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.nf') as tmp:
        tmp.write(content)
        tmp.flush()
        
        # Create module from file
        module_result = Module.from_file(tmp.name)
        assert module_result.error is not None
        assert module_result.error.likely_rt_bug == False
        

def test_parse_modules():
    # Create a temporary directory
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a valid module file
        valid_content = """
        process VALID_PROCESS {
            input:
                path input_file
            
            output:
                path "output.txt"
            
            script:
            '''
            echo "processing $input_file" > output.txt
            '''
        }
        """
        
        # Create an invalid module file (using undefined MyClass)
        invalid_content = """
        class MyClass {
            // Invalid - missing return type and modifiers
            myMethod() {
                println "Hello"
            }
        }
        """
        
        # Write the files
        with open(os.path.join(tmpdir, "valid.nf"), "w") as f:
            f.write(valid_content)
        with open(os.path.join(tmpdir, "invalid.nf"), "w") as f:
            f.write(invalid_content)
            
        # Track progress
        progress_calls = []
        def progress_callback(current, total):
            progress_calls.append((current, total))
            
        # Process the modules
        results: List[ModuleResult] = parse_modules(tmpdir, progress_callback)
        
        # Verify progress tracking
        assert len(progress_calls) == 2  # Should be called twice
        assert progress_calls[-1] == (2, 2)  # Final call should be (2, 2)
        
        # Check results
        assert len(results) == 2
        
        # Count successes and failures
        successes = [r for r in results if r.module is not None]
        failures = [r for r in results if r.error is not None]
        
        assert len(successes) == 1
        assert len(failures) == 1
        
        # Verify the valid module
        valid_result = next(r for r in results if r.module is not None)
        assert "VALID_PROCESS" in valid_result.module.processes[0].name
        
        # Verify the invalid module
        invalid_result = next(r for r in results if r.error is not None)
        assert "myMethod" in invalid_result.error.error  # Error should mention MyClass
        