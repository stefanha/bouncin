"""Live update module reloading support"""

import sys

__all__ = ['run', 'Reload']

# Note that the reloader module never reloads itself.  It may be possible to
# reload most of this module but I prefer to keep the code short and simple.

def _get_current_modules():
    return set(sys.modules.keys())

def _unload_modules(baseline_modules):
    curmods = _get_current_modules()
    for modname in curmods.difference(baseline_modules):
        del sys.modules[modname]

class Reload(Exception):
    """Raise this exception to reload all modules"""
    pass

def run(start_modname):
    """Import a main module, reload it and all modules if Reload is raised"""

    # Snapshot currently loaded modules, these are considered default runtime
    # modules and will never be reloaded.  Not all Python modules can be
    # unloaded so it is necessary to keep this list.
    baseline_modules = _get_current_modules()

    while True:
        try:
            __import__(start_modname)
            return
        except Reload:
            _unload_modules(baseline_modules)
