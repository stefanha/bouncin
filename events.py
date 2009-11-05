"""Event notification system"""

__all__ = ['add_handler', 'notify']

# There is an event handler list for each named event
_handlers = {}

def add_handler(event_name, handler_fn):
    """Register a handler function for a named event"""
    fns = _handlers.get(event_name, [])
    if handler_fn not in fns:
        fns.append(handler_fn)
    _handlers[event_name] = fns

def notify(event_name, *args):
    """Call each handler function for an event"""
    for handler_fn in _handlers.get(event_name, []):
        handler_fn(*args)
