package core

import (
	"os";
	"container/list";
	"net";
	"log";
	"irc";
	"runloop";
)

type Network struct {
	name		string;
	server		*server;
	clients		*list.List;
	listen		*listenConn;
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
	network = &Network{name: name, clients: list.New(), listen: l};
	network.server = newServer(serverConn, network);
	return network;
}

func (network *Network) addClient(conn net.Conn) {
	// TODO error handler on client for disconnect
	client := newClient(conn, network);
	network.clients.PushBack(client);
	log.Stderrf("client connected from %s\n", conn.RemoteAddr());
}

func (network *Network) SendToServer(msg *irc.Message) {
	// TODO network may be down, but is this the way to handle it?
	if network.server != nil {
		network.server.Send(msg)
	}
}

func (network *Network) SendToClients(msg *irc.Message) {
	for c := range network.clients.Iter() {
		c.(*client).Send(msg)
	}
}

var networks = make(map[string] *Network);

func AddNetwork(name string, server net.Conn, listen net.Listener) *Network {
	// TODO what if network already exists?
	network := newNetwork(name, server, listen);
	networks[name] = network;
	return network;
}
