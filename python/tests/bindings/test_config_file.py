import pytest
import tempfile
from pathlib import Path
from reftrace import ConfigFile

def test_parse_config_conditional():
    """Test parsing of conditional process blocks"""
    content = """
if (!params.skip_bbsplit && params.bbsplit_fasta_list) {
    process {
        withName: '.*:PREPARE_GENOME:BBMAP_BBSPLIT' {
            ext.args   = 'build=1'
            publishDir = [
                path: { "${params.outdir}/genome/index" },
                mode: params.publish_dir_mode,
                saveAs: { filename -> filename.equals('versions.yml') ? null : filename },
                enabled: params.save_reference
            ]
        }
    }
}
"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.config') as tmp:
        tmp.write(content)
        tmp.flush()
        
        config = ConfigFile.from_file(tmp.name)
        
        # Should have one process scope
        assert len(config.process_scopes) == 1
        
        # Check process scope content
        scope = config.process_scopes[0]
        assert len(scope.named_scopes) == 1
        
        # Check named scope
        named_scope = scope.named_scopes[0]
        assert named_scope.name == ".*:PREPARE_GENOME:BBMAP_BBSPLIT"
        
        # Check directives
        assert len(named_scope.directives) == 2
        args_directive = next(d for d in named_scope.directives if d.name == "ext.args")
        publish_directive = next(d for d in named_scope.directives if d.name == "publishDir")
        
        # Check directive values according to proto definition
        assert len(args_directive.value.params) == 0  # No params in static value
        assert not args_directive.value.in_closure  # Not in closure
        
        # Check publishDir options
        assert len(publish_directive.options) == 4  # path, mode, saveAs, enabled
        path_option = next(o for o in publish_directive.options if o.name == "path")
        assert path_option.value.in_closure  # path is in closure
        assert "outdir" in path_option.value.params  # should contain outdir param

def test_parse_config_duplicate_arg():
    """Test parsing of complex argument expressions"""
    content = """
        process {
            withName: '.*:FASTQ_FASTQC_UMITOOLS_TRIMGALORE:TRIMGALORE' {
                ext.args   = {
                    [
                        "--fastqc_args '-t ${task.cpus}'",
                        params.extra_trimgalore_args ? params.extra_trimgalore_args.split("\\s(?=--)") : ''
                    ].flatten().unique(false).join(' ').trim()
                }
        }
    }"""
    
    with tempfile.NamedTemporaryFile(mode='w', suffix='.config') as tmp:
        tmp.write(content)
        tmp.flush()
        
        config = ConfigFile.from_file(tmp.name)
        
        # Basic structure validation
        assert len(config.process_scopes) == 1
        process_scope = config.process_scopes[0]
        assert len(process_scope.named_scopes) == 1
        
        # Named scope validation
        named_scope = process_scope.named_scopes[0]
        assert named_scope.name == ".*:FASTQ_FASTQC_UMITOOLS_TRIMGALORE:TRIMGALORE"
        
        # Directive validation
        assert len(named_scope.directives) == 1
        directive = named_scope.directives[0]
        assert directive.name == "ext.args"
        
        # Value validation based on proto definition
        value = directive.value
        assert len(value.params) == 1
        assert value.params[0] == "extra_trimgalore_args"
        assert value.in_closure == True