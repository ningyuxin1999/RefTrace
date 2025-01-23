import tempfile
from reftrace import Module, parse_modules, ParseError
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
        module = Module.from_file(tmp.name)
        assert not isinstance(module, ParseError)
        
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
        assert isinstance(module_result, ParseError)
        assert module_result.likely_rt_bug == False
        

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
        result = parse_modules(tmpdir, progress_callback)
        
        # Verify progress tracking
        assert len(progress_calls) == 2  # Should be called twice
        assert progress_calls[-1] == (2, 2)  # Final call should be (2, 2)
        
        # Check results
        assert len(result.results) == 1
        assert len(result.errors) == 1
        
        # Verify the valid module
        valid_result = result.results[0]
        assert "VALID_PROCESS" in valid_result.processes[0].name
        
        # Verify the invalid module
        invalid_result: ParseError = result.errors[0]
        assert "myMethod" in invalid_result.error  # Error should mention MyClass
        