import pytest
from click.testing import CliRunner
import json
import tempfile
from pathlib import Path
from reftrace.cli.info import show_modules_info, rdeps

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

def test_rdeps_basic():
    """Test basic reverse dependencies output"""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a simple workflow with dependencies
        main_content = """
        include { FOO } from './modules/foo/main'
        include { BAR } from './modules/bar/main'
        
        workflow {
            FOO()
            BAR(FOO.out)
        }
        """
        (Path(tmpdir) / "main.nf").write_text(main_content)
        
        # Create modules directory and module subdirectories
        modules_dir = Path(tmpdir) / "modules"
        modules_dir.mkdir()
        
        # Create foo module
        foo_dir = modules_dir / "foo"
        foo_dir.mkdir()
        (foo_dir / "main.nf").write_text("""
        process FOO {
            output: path "foo.txt"
            script: "echo foo > foo.txt"
        }
        """)
        
        # Create bar module
        bar_dir = modules_dir / "bar"
        bar_dir.mkdir()
        (bar_dir / "main.nf").write_text("""
        process BAR {
            input: path x
            script: "echo bar"
        }
        """)
        
        # Run command
        runner = CliRunner(mix_stderr=False)
        result = runner.invoke(rdeps, ['--directory', tmpdir])
        
        # Check command executed successfully
        assert result.exit_code == 0
        
        # Parse JSON output
        output = json.loads(result.stdout)
        
        # Verify structure and content
        assert isinstance(output, list)
        
        # Find foo module entry
        foo_entry = next(m for m in output if m["path"].endswith("modules/foo/main.nf"))
        assert len(foo_entry["direct_rdeps"]) == 1
        assert "main.nf" in foo_entry["direct_rdeps"][0]
        assert foo_entry["transitive_rdeps"] == []
        
        # Find bar module entry
        bar_entry = next(m for m in output if m["path"].endswith("modules/bar/main.nf"))
        assert len(bar_entry["direct_rdeps"]) == 1
        assert "main.nf" in bar_entry["direct_rdeps"][0]
        assert bar_entry["transitive_rdeps"] == []

def test_rdeps_isolated_nodes():
    """Test detection of isolated nodes"""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create an isolated module
        (Path(tmpdir) / "isolated.nf").write_text("""
        process ISOLATED {
            script: "echo isolated"
        }
        """)
        
        # Run command with --isolated flag
        runner = CliRunner(mix_stderr=False)
        result = runner.invoke(rdeps, ['--directory', tmpdir, '--isolated'])
        
        # Should exit with code 1 due to isolated node
        assert result.exit_code == 1
        assert "Warning: Found isolated nodes" in result.stderr
        assert "isolated.nf" in result.stderr

def test_rdeps_pretty_vs_compact():
    """Test pretty vs compact JSON output"""
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create a simple module
        (Path(tmpdir) / "test.nf").write_text("""
        process TEST {
            script: "echo test"
        }
        """)
        
        # Test pretty output
        runner = CliRunner(mix_stderr=False)
        pretty_result = runner.invoke(rdeps, ['--directory', tmpdir, '--pretty'])
        
        # Test compact output
        compact_result = runner.invoke(rdeps, ['--directory', tmpdir, '--compact'])
        
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
