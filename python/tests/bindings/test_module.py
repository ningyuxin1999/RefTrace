import pytest
import tempfile
from pathlib import Path
from reftrace import Module, ContainerDirective

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
        
        # Should have one process
        assert len(module.processes) == 1
        process = module.processes[0]
        
        # Check process name
        assert process.name == "CAT_FASTQ"

        containers = process.get_directives("container")
        
        # Check container directive
        containers = [c for c in containers if c.format == ContainerDirective.Format.TERNARY]
        assert len(containers) == 1
        container = containers[0]
        
        # Verify container properties
        assert container.condition == "((workflow.containerEngine == singularity) && !(task.ext.singularity_pull_docker_container))"
        assert container.true_name == "https://depot.galaxyproject.org/singularity/ubuntu:20.04"
        assert container.false_name == "nf-core/ubuntu:20.04"
        
        # Add label directive testing
        labels = process.get_directives("label")
        assert len(labels) == 1
        label = labels[0]
        assert label.label == "process_single"