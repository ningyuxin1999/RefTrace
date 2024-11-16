import pytest
import tempfile
from pathlib import Path
from reftrace.cli.main import run_lint

def create_test_file(content: str) -> tuple[str, Path]:
    """Helper function to create a test file and return its path"""
    tmpdir = tempfile.mkdtemp()
    nf_path = Path(tmpdir) / "workflow.nf"
    nf_path.write_text(content)
    
    rules_path = Path(__file__).parent.parent / "fixtures" / "rules.py"

    return tmpdir, rules_path

def test_container_with_space():
    content = """
process FOO {
    container 'ubuntu latest'  // Invalid container name with space
    cpus 4
    memory '8 GB'

    script:
    \"\"\"
    echo "test"
    \"\"\"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 2
    assert any("container name 'ubuntu latest' contains spaces, which is not allowed" in err.error for err in result.errors)
    assert any("unknown container type" in err.error for err in result.errors)
def test_container_without_space():
    content = """
process FOO {
    label "process_single"
    container 'ubuntu:latest'  // Valid container name without space
    cpus 4
    memory '8 GB'

    script:
    \"\"\"
    echo "test"
    \"\"\"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 0
    assert len(result.warnings) == 0

def test_container_with_mixed_syntax():
    content = """
process FOO {
    label "process_single"
    container "biocontainers/ubuntu https://containers/something"
    script:
    \"\"\"
    echo "test"
    \"\"\"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.warnings) == 1
    assert "Docker and Singularity containers specified on the same line" in result.warnings[0].warning

def test_valid_docker_container_with_tag():
    content = """
process FOO {
    container "biocontainers/fastqc:0.11.9"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 0

def test_valid_singularity_container_with_tag():
    content = """
process FOO {
    container "https://depot.galaxyproject.org/singularity/fastqc:0.11.9--0"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 0

def test_docker_container_missing_tag():
    content = """
process FOO {
    container "biocontainers/fastqc"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 1
    assert "docker container 'biocontainers/fastqc' must specify a tag" in result.errors[0].error

def test_singularity_container_missing_tag():
    content = """
process FOO {
    container "https://depot.galaxyproject.org/singularity/fastqc"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 1
    assert "singularity container" in result.errors[0].error
    assert "must specify a tag" in result.errors[0].error

def test_quay_io_container_without_tag():
    content = """
process FOO {
    container "quay.io/biocontainers/fastqc"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 2
    assert any("must specify a tag" in err.error for err in result.errors)
    assert any("please use 'organization/container:tag' format" in err.error for err in result.errors)

def test_multiple_containers_with_ternary():
    content = """
process FOO {
    container "${workflow.containerEngine == 'singularity' ? 
        'https://depot.galaxyproject.org/singularity/fastqc' : 
        'biocontainers/fastqc'}"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 2
    assert any("singularity container" in err.error and "must specify a tag" in err.error for err in result.errors)
    assert any("docker container" in err.error and "must specify a tag" in err.error for err in result.errors)

def test_unknown_container_type():
    content = '''
process FOO {
    container "just-a-name"  // Invalid: no URL prefix or org/image format
    input:
    path x

    output:
    path "y"

    script:
    """
    echo "test"
    """
}
'''
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 1
    assert "unknown container type" in result.errors[0].error

def test_invalid_singularity_container_url():
    content = """
process FOO {
    container "https://[invalid-url"  // Invalid: malformed URL
    script:
    \"\"\"
    echo "test"
    \"\"\"
}
"""
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 1
    assert "invalid singularity container URL" in result.errors[0].error

def test_invalid_docker_tag_format():
    content = '''
process FOO {
    container "biocontainers/fastqc:1.0!"  // Invalid characters in tag
    script:
    """
    echo "test"
    """
}
'''
    tmpdir, rules_path = create_test_file(content)
    results = run_lint(tmpdir, str(rules_path))
    
    assert len(results) == 1
    result = results[0]
    assert len(result.errors) == 1
    assert "invalid docker tag format" in result.errors[0].error