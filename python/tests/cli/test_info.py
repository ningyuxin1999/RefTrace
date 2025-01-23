import pytest
from click.testing import CliRunner
import json
import tempfile
from pathlib import Path
from reftrace.cli.info import show_modules_info

def test_show_modules_info_basic():
    """Test basic module info output"""
    # Create temporary test files
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a simple valid module file
        content = """
        process FOO {
            input:
                path x
            output:
                path "y"
            script:
            '''
            echo "test"
            '''
        }
        """
        module_path = Path(tmpdir) / "workflow.nf"
        module_path.write_text(content)
        
        # Run command with stdout captured but stderr ignored
        runner = CliRunner(mix_stderr=False)
        result = runner.invoke(show_modules_info, ['--directory', tmpdir])
        
        # Check command executed successfully
        assert result.exit_code == 0
        
        # Parse JSON output from stdout
        output = json.loads(result.stdout)
        
        # Check structure
        assert "modules" in output
        assert "resolved_includes" in output
        assert "unresolved_includes" in output
        
        # Check module content
        assert len(output["modules"]) == 1
        module = output["modules"][0]
        assert "path" in module
        assert str(module_path) in module["path"]
        assert "processes" in module
        assert len(module["processes"]) == 1
        assert module["processes"][0]["name"] == "FOO"

def test_show_modules_info_paths_only():
    """Test --paths flag that shows only module paths"""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create test files
        content1 = "process FOO { script: 'echo test' }"
        content2 = "process BAR { script: 'echo test' }"
        
        path1 = Path(tmpdir) / "workflow1.nf"
        path2 = Path(tmpdir) / "workflow2.nf"
        path1.write_text(content1)
        path2.write_text(content2)
        
        # Run command with --paths flag
        runner = CliRunner(mix_stderr=False)
        result = runner.invoke(show_modules_info, ['--directory', tmpdir, '--paths'])
        
        # Check command executed successfully
        assert result.exit_code == 0
        
        # Parse JSON output
        paths = json.loads(result.stdout)
        
        # Check it's a list of paths
        assert isinstance(paths, list)
        assert len(paths) == 2
        assert all(isinstance(p, str) for p in paths)
        assert all(p.endswith('.nf') for p in paths)

def test_show_modules_info_invalid_file():
    """Test handling of invalid Nextflow files"""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create invalid module file
        # TODO: try with a comment next to invalid_directive.
        # not sure why that doesn't cause an error.
        content = """
        process FOO {
            invalid_directive
        }
        """
        module_path = Path(tmpdir) / "invalid.nf"
        module_path.write_text(content)
        
        # Run command
        runner = CliRunner(mix_stderr=False)
        result = runner.invoke(show_modules_info, ['--directory', tmpdir])
        
        # Command should complete but with error message
        assert result.exit_code == 0
        
        # Parse JSON output from stdout
        output = json.loads(result.stdout)
        
        # Should have empty structure
        assert output == {
            "modules": [],
            "resolved_includes": [],
            "unresolved_includes": []
        }

def test_show_modules_info_pretty_vs_compact():
    """Test pretty vs compact JSON output"""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a simple module file
        content = "process FOO { script: 'echo test' }"
        module_path = Path(tmpdir) / "workflow.nf"
        module_path.write_text(content)
        
        # Test pretty output
        runner = CliRunner(mix_stderr=False)
        pretty_result = runner.invoke(show_modules_info, ['--directory', tmpdir, '--pretty'])
        
        # Test compact output
        compact_result = runner.invoke(show_modules_info, ['--directory', tmpdir, '--compact'])
        
        # Both should succeed
        assert pretty_result.exit_code == 0
        assert compact_result.exit_code == 0
        
        # Pretty should have newlines and indentation
        assert '\n' in pretty_result.stdout
        assert '  ' in pretty_result.stdout
        
        # Compact should be single line
        assert len(compact_result.stdout.strip().split('\n')) == 1
        
        # Both should parse to equivalent JSON
        assert json.loads(pretty_result.stdout) == json.loads(compact_result.stdout)