package network

import (
	"os";
	"container/list";
	"net";
	"log";
	"runloop";
	"server";
	"client";
)


type acceptFunc	func(net.Conn)
type errorFunc	func(os.Error)

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


type Network struct {
	name		string;
	server		*server.Server;
	clients		*list.List;
	listen		*listenConn;
}

var networks = make(map[string] *Network);

func Add(name string, server net.Conn, listen net.Listener) *Network {
	// TODO what if network already exists?
	network := newNetwork(name, server, listen);
	networks[name] = network;
	return network;
}

func newNetwork(name string, serverConn net.Conn, listen net.Listener) *Network {
	var network *Network;

	accept := func(conn net.Conn) {
		runloop.CallLater(func() {
			network.addClient(conn)
		})
	};

	error := func(err os.Error) {
		// TODO listener failed
	};

	l := newListenConn(listen, accept, error);
	network = &Network{name, server.New(serverConn), list.New(), l};
	return network;
}

func (network *Network) addClient(conn net.Conn) {
	// TODO error handler on client for disconnect
	client := client.New(conn);
	network.clients.PushBack(client);
	log.Stderrf("client connected from %s\n", conn.RemoteAddr());
}
