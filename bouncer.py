"""Main bouncer code"""

import reloader
import events
import runloop

try:
    runloop.run()
except reloader.Reload:
    # Notify so modules can stash state before being reloaded
    events.notify('pre-reload')
    raise
