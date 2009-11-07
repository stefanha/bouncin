"""IRC protocol

See RFC 2812 "Internet Relay Chat: Client Protocol"."""

import re

__all__ = ['Message', 'parse', 'format']

MESSAGE_RE = re.compile(r'(?::([^ ]+) )?([^ ]+)(.*)')

class Message(object):
    """An IRC message"""

    def __init__(self, command, *params, **kwargs):
        self.command = command
        self.params = params
        self.prefix = kwargs.get('prefix', '')

    def __str__(self):
        fields = []

        if self.prefix:
            fields.append(':' + self.prefix)

        fields.append(self.command)

        if self.params and ' ' in self.params[-1]:
            fields += self.params[:-1]
            fields.append(':' + self.params[-1])
        else:
            fields += self.params

        return ' '.join(fields)

def parse(line, from_client=True):
    """Break up an IRC message line into a Message object"""

    m = MESSAGE_RE.match(line)
    if not m:
        raise ValueError

    prefix, command, params = m.groups()

    # Discard prefix if sent from a client
    if not prefix or from_client:
        prefix = ''

    # Convert parameters into a list of strings
    trailing = []
    if ' :' in params:
        params, trailing = params.rsplit(' :', 1)
        trailing = [trailing]
    params = params.lstrip().split() + trailing

    return Message(command, prefix=prefix, *params)

def format(command, *args, **kwargs):
    """Build an IRC message line from arguments"""
    return str(Message(command, *args, **kwargs))
