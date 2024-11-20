# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

import os
import sys
sys.path.insert(0, os.path.abspath('../python'))

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

# Autodoc settings
autodoc_default_options = {
    'members': True,
    'imported-members': True
}

napoleon_use_rtype = False

autodoc_typehints = 'signature'

# don't show signatures for classes
autodoc_class_signature = 'separated'

# Mock the C extension
autodoc_mock_imports = ['libreftrace']

markdown_anchor_signatures= True

autodocsumm_section_sorter = False

def skip_init(app, what, name, obj, skip, options):
    if name == "__init__" and what == "class" and obj.__qualname__ != "Module.__init__":
        return True
    return skip

def setup(app):
    app.connect('autodoc-skip-member', skip_init)

