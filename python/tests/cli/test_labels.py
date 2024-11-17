import pytest
import tempfile
from pathlib import Path
from reftrace.cli.main import run_lint
from tests.utils import create_test_file

def test_no_labels():
    content = """
process FOO {
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
    assert len(result.warnings) == 2
    assert any("has no labels" in w.warning for w in result.warnings)
    assert any("has no standard label" in w.warning for w in result.warnings)

def test_standard_label():
    content = """
process FOO {
    label 'process_single'
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
    assert len(result.warnings) == 0

def test_conflicting_labels():
    content = """
process FOO {
    label 'process_single'
    label 'process_high'
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
    assert "has conflicting labels" in result.warnings[0].warning
    assert "process_single" in result.warnings[0].warning
    assert "process_high" in result.warnings[0].warning

def test_non_standard_label():
    content = """
process FOO {
    label 'custom_label'
    label 'another_label'
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
    assert len(result.warnings) == 2  # Should have no_standard_label and non_standard_label warnings
    assert any("has non-standard labels: ['custom_label', 'another_label']" in w.warning for w in result.warnings)
    assert any("has no standard label" in w.warning for w in result.warnings)

def test_duplicate_labels():
    content = """
process FOO {
    label 'process_single'
    label 'process_single'
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
    assert len(result.warnings) == 2
    assert any("has duplicate label 'process_single' (2 times)" in w.warning for w in result.warnings)
    assert any("has conflicting labels" in w.warning for w in result.warnings)

def test_non_alphanumeric_label():
    content = """
process FOO {
    label 'process_single'
    label 'label@123'
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
    assert len(result.warnings) == 2
    assert any("contains non-alphanumeric characters" in w.warning for w in result.warnings)
    assert any("has non-standard labels" in w.warning for w in result.warnings)