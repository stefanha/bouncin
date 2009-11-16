package main

import (
	"container/list";
	"net";
	"log";
	"irc";
)

type Network struct {
	name	string;
	server	*irc.Conn;
	listen	net.Listener;
	clients	*list.List;
}

func NewNetwork(name string, server net.Conn, listen net.Listener) *Network {
	if listen == nil {
		return nil
	}
	network := &Network{name, irc.NewConn(server, nil), listen, list.New()};
	go network.run();
	return network;
}

func (network *Network) acceptor(accepted chan net.Conn) {
	for {
		conn, err := network.listen.Accept();
		if err != nil {
			log.Stderrf("accept failed: %s\n", err);
			// TODO handle error
			return
		}
		accepted <- conn;
	}
}

func (network *Network) run() {
	accepted := make(chan net.Conn, 0);
	go network.acceptor(accepted);
	log.Stderrf("network %s listening on %s\n", network.name, network.listen.Addr());
	for {
		select {
		case client := <-accepted:
			network.addClient(client)
		}
	}
}

func (network *Network) addClient(conn net.Conn) {
	client := irc.NewConn(conn, nil);
	network.clients.PushBack(client);
	log.Stderrf("client connected from %s\n", conn.RemoteAddr());
}
