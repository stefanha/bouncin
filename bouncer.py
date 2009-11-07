"""Main bouncer code"""

import socket
import asyncore
import asynchat
import logging
import reloader
import events
import runloop

__all__ = []

class LineBasedSocket(asynchat.async_chat):
    def __init__(self, conn=None):
        self._ibuffer = []
        asynchat.async_chat.__init__(self, conn)
        self.set_terminator('\r\n')

    def collect_incoming_data(self, data):
        self._ibuffer.append(data)

    def found_terminator(self):
        data = ''.join(self._ibuffer)
        self._ibuffer = []
        events.notify('socket.line-recv', self, data)
        logging.debug('%s: line-recv: %s' % (__name__, data))

    def send_line(self, data):
        """Send a line via this socket

        A newline is automatically appended."""

        logging.debug('%s: line-send: %s' % (__name__, data))
        data += '\r\n'
        events.notify('socket.line-send', self, data)
        self.push(data)

class ServerSocket(asyncore.dispatcher):
    def __init__(self, addr):
        asyncore.dispatcher.__init__(self)
        self.create_socket(socket.AF_INET, socket.SOCK_STREAM)
        self.set_reuse_addr()
        self.bind(addr)
        self.listen(1)
        logging.info('%s: listening on %s' % (__name__, str(addr)))

    def handle_accept(self):
        conn, addr = self.accept()
        logging.info('%s: accept from %s' % (__name__, str(addr)))
        events.notify('socket.accept', self, conn, addr)

logging.getLogger().setLevel(logging.DEBUG)

ServerSocket(('0.0.0.0', 10000)) # TODO

def handle_accept(server_socket, conn, addr):
    LineBasedSocket(conn)
events.add_handler('socket.accept', handle_accept)

try:
    runloop.run()
except reloader.Reload:
    # Notify so modules can stash state before being reloaded
    events.notify('pre-reload')

    logging.info('reloading (live code update)...')
    raise
