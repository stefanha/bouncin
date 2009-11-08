"""Main bouncer code"""

import socket
import asyncore
import asynchat
import logging
import reloader
import events
import runloop
import irc
import config

__all__ = []

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
        logging.debug('%s: line-recv: %s' % (__name__, data))
        events.notify('socket.line-recv', self, data)

    def send_line(self, data):
        """Send a line via this socket

        A newline is automatically appended."""

        logging.debug('%s: line-send: %s' % (__name__, data))
        data += '\r\n'
        events.notify('socket.line-send', self, data)
        self.push(data)

class Network(object):
    def __init__(self, name, server_addr, listen_addr):
        self.name = name
        # TODO connect to server
        self.listen = ServerSocket(listen_addr)
        self.clients = []
        self.nick = 'nicknotset'

        events.add_handler('socket.accept', self._on_accept)
        events.add_handler('socket.line-recv', self._on_line_recv)
        events.add_handler('irc.recv', self._on_irc_recv)

    def _on_accept(self, server_sock, conn, addr):
        if server_sock is not self.listen:
            return
        self.clients.append(LineBasedSocket(conn))

    def _on_line_recv(self, sock, line):
        # TODO process server messages
        if sock not in self.clients:
            return

        try:
            msg = irc.parse(line)
        except ValueError:
            logging.warn('%s: invalid IRC line: %s' % (__name__, line))
            return

        events.notify('irc.recv', sock, msg)

    def _on_irc_recv(self, sock, msg):
        if sock not in self.clients:
            return

        if msg.command == 'NICK':
            self.nick = msg.params[0]
        if msg.command == 'USER':
            # TODO complete fields
            sock.send_line(irc.format(irc.RPL_WELCOME, self.nick, 'Welcome to the Internet Relay Network %s!<user>@<host>' % self.nick))
            sock.send_line(irc.format(irc.RPL_YOURHOST, self.nick, 'Your host is <servername>, running version <ver>'))
            sock.send_line(irc.format(irc.RPL_CREATED, self.nick, 'This server was created <date>'))
            sock.send_line(irc.format(irc.RPL_MYINFO, self.nick, '<servername> <version> <available user modes> <available channel modes>'))

def hostname_to_addr(hostname, port=6667):
    if ':' in hostname:
        host, port = hostname.split(':')
        try:
            port = int(port)
        except ValueError:
            host = hostname
    return host, port

logging.getLogger().setLevel(logging.DEBUG)

networks = []
cfg = config.cfg
for section in cfg.sections():
    # Check if this config section defines a network
    if not cfg.has_option(section, 'listen-addr'):
        continue

    name = section
    server_addr = None # TODO
    listen_addr = hostname_to_addr(cfg.get('listen-addr'))

    networks.append(Network(name, server_addr, listen_addr))

# Bring up configuration "network" if no networks are defined
if not networks:
    host = '0.0.0.0'
    port = cfg.getint('bouncer', 'port_range_start')
    networks.append(Network('bouncin', None, (host, port)))

    print 'Please connect to %s:%s with an IRC client to get started' % (host, port)

try:
    runloop.run()
except reloader.Reload:
    # Notify so modules can stash state before being reloaded
    events.notify('pre-reload')

    logging.info('reloading (live code update)...')
    raise
