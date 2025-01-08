import os
import sys
import ctypes
from ctypes import c_char_p, c_void_p

def load_library():
    _lib_dir = os.path.dirname(os.path.abspath(__file__))

    # Platform-specific library names
    if sys.platform == "darwin":
        lib_name = "libreftrace.dylib"
    elif sys.platform == "win32":
        lib_name = "libreftrace.dll"
    else:  # Linux and others
        lib_name = "libreftrace.so"

    _lib_path = os.path.join(_lib_dir, lib_name)
    
    if not os.path.exists(_lib_path):
        raise RuntimeError(f"Library not found at {_lib_path}")
        
    try:
        return ctypes.CDLL(_lib_path)
    except OSError as e:
        raise RuntimeError(f"Failed to load library: {e}") from e

# Load the shared library
_lib = load_library()

_lib.Module_New.argtypes = [c_char_p]
_lib.Module_New.restype = c_void_p

_lib.Module_Free.argtypes = [c_void_p]
_lib.Module_Free.restype = None

_lib.ConfigFile_New.argtypes = [c_char_p]
_lib.ConfigFile_New.restype = c_void_p

_lib.ConfigFile_Free.argtypes = [c_void_p]
_lib.ConfigFile_Free.restype = None

_lib.Parse_Modules.argtypes = [c_char_p, c_void_p]
_lib.Parse_Modules.restype = c_void_p
