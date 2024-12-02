# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

import os
import sys
sys.path.insert(0, os.path.abspath('../../python'))

project = 'RefTrace'
copyright = '2024'
author = 'Andrew Stiles'

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = [
    'sphinx.ext.autodoc',
    'sphinx.ext.autosummary',
    'sphinx.ext.napoleon',
    'autodocsumm'
]


napoleon_use_rtype = False

autodoc_typehints = 'signature'

# don't show signatures for classes
autodoc_class_signature = 'separated'

markdown_anchor_signatures= True

autodocsumm_section_sorter = False

# Add autosummary settings
# autosummary_generate = True  # Generate stub pages
add_module_names = False     # Don't prefix with module names

# Mock the protobuf modules and specific classes
autodoc_mock_imports = [
    'libreftrace',  # The C extension
    'proto.config_file_pb2',
    'proto.module_pb2',
    'proto.config_file_pb2.ConfigFile',
    'proto.config_file_pb2.ProcessScope',
    'proto.module_pb2.Module',
    'proto.module_pb2.Process',
    'proto'
]

# Important for handling aliases
autoclass_content = 'class'  # Only use class docstring, not __init__
autodoc_inherit_docstrings = True
autodoc_member_order = 'bysource'

# Don't be so strict about references
nitpicky = False
