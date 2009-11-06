"""Main loop for an asynchronous system"""

import time
import sched
import asyncore

__all__ = ['Runloop', 'StopRunloop']

class StopRunloop(Exception):
    pass

def _delayfunc(duration):
    """Process I/O until duration expires
    
    If the scheduler has no work to do, then process I/O until there is
    more work."""

    # Stop scheduler if all I/O is closed
    if not asyncore.socket_map:
        raise StopRunloop()

    if duration:
        # Process I/O for one iteration or timeout when duration expires
        asyncore.loop(duration, use_poll=True, count=1)
    else:
        # Process I/O until the scheduler has work to do again
        while _scheduler.empty():
            asyncore.loop(use_poll=True, count=1)

_scheduler = sched.scheduler(time.time, _delayfunc)

def add_oneshot_timer(duration, handler_fn, args):
    """Schedule a timer to fire duration seconds into the future"""
    _scheduler.enter(duration, 1, handler_fn, args)

def run():
    """Run forever processing I/O and timers
    
    Return when all I/O is closed."""
    try:
        _scheduler.run()
    except StopRunloop:
        pass
