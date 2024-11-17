import tempfile
from pathlib import Path
from importlib.resources import files

def create_test_file(content: str) -> tuple[str, Path]:
    """Helper function to create a test file and return its path.
    
    Args:
        content: The content to write to the test file
        
    Returns:
        tuple: (temporary directory path, rules file path)
    """
    tmpdir = tempfile.mkdtemp()
    nf_path = Path(tmpdir) / "workflow.nf"
    nf_path.write_text(content)
    
    # Get rules.py from package resources
    rules_path = files('reftrace.fixtures') / 'rules.py'

    return tmpdir, rules_path