package core

import (
	"net";
	"os";
	"log";
)

type acceptFunc	func(net.Conn)
type errorFunc	func(os.Error)

// A listen socket.
type listenConn struct {
	listen	net.Listener;
	accept	acceptFunc;
	error	errorFunc;
}

func newListenConn(listen net.Listener, accept acceptFunc, error errorFunc) *listenConn {
	l := &listenConn{listen, accept, error};
	go l.run();
	return l;
}

func (l *listenConn) run() {
	log.Stderrf("listening on %s\n", l.listen.Addr());
	// TODO a way to shut this down
	for {
		conn, err := l.listen.Accept();
		if err != nil {
			log.Stderrf("accept failed: %s\n", err);
			l.listen.Close();
			if l.error != nil {
				l.error(err);
			}
			return
		}
		l.accept(conn);
	}
}

func (l *listenConn) Addr() net.Addr {
	return l.listen.Addr()
}
