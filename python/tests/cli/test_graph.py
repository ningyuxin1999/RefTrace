import pytest
import tempfile
from pathlib import Path
import os
from reftrace.cli.graph import graph
from click.testing import CliRunner
from tests.utils import create_test_file

def test_graph_generation():
    """Test basic graph generation with a simple workflow"""
    content = """
    include { FOO } from './modules/foo'
    include { BAR } from './modules/bar'
    
    workflow {
        FOO()
        BAR()
    }
    """
    
    # Create main workflow file
    tmpdir = tempfile.mkdtemp()
    main_path = Path(tmpdir) / "workflow.nf"
    main_path.write_text(content)
    
    # Create module files
    modules_dir = Path(tmpdir) / "modules"
    modules_dir.mkdir()
    
    # Create foo.nf
    foo_content = """
    process FOO {
        output:
        path "foo.txt"
        
        script:
        '''
        echo "foo" > foo.txt
        '''
    }
    """
    (modules_dir / "foo.nf").write_text(foo_content)
    
    # Create bar.nf
    bar_content = """
    process BAR {
        output:
        path "bar.txt"
        
        script:
        '''
        echo "bar" > bar.txt
        '''
    }
    """
    (modules_dir / "bar.nf").write_text(bar_content)
    
    # Run graph command
    runner = CliRunner()
    result = runner.invoke(graph, ['--directory', tmpdir])
    
    # Check command executed successfully
    assert result.exit_code == 0
    
    # Check output file was created
    graph_path = Path("graph.png")
    assert graph_path.exists()
    
    # Clean up
    graph_path.unlink()

def test_graph_with_errors():
    """Test graph generation with invalid Nextflow files"""
    content = """
    // Invalid Nextflow syntax
    process FOO {
        invalid_directive
    }
    """
    
    tmpdir = tempfile.mkdtemp()
    main_path = Path(tmpdir) / "workflow.nf"
    main_path.write_text(content)
    
    runner = CliRunner()
    result = runner.invoke(graph, ['--directory', tmpdir])
    
    # Should exit with error
    assert result.exit_code == 1
    assert "Failed to parse" in result.output

def test_graph_empty_directory():
    """Test graph generation with empty directory"""
    tmpdir = tempfile.mkdtemp()
    
    runner = CliRunner()
    result = runner.invoke(graph, ['--directory', tmpdir])
    
    assert result.exit_code == 1
    assert not Path("graph.png").exists()

def test_graph_nonexistent_directory():
    """Test graph generation with non-existent directory"""
    runner = CliRunner()
    result = runner.invoke(graph, ['--directory', '/nonexistent/path'])
    
    # Should fail with error
    assert result.exit_code == 2
    assert "does not exist" in result.output.lower()