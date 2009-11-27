package core

import (
	"os";
	"net";
	"log";
	"irc";
	"runloop";
	"events";
)

// A server implements a connection to an IRC server.
type server struct {
	conn	*irc.Conn;
	network	*Network;
}

// newServer returns a new server for a given network connection.
func newServer(conn net.Conn, network *Network) *server {
	var s *server;

	if conn == nil {
		return nil
	}

	recvFunc := func(msg *irc.Message) {
		runloop.CallLater(func() { s.recvFunc(msg) })
	};
	errorFunc := func(err os.Error) {
		runloop.CallLater(func() { s.errorFunc(err) })
	};

	s = &server{irc.NewConn(conn, recvFunc, errorFunc), network};
	s.conn.Send(&irc.Message{Command: "NICK", Params: []string{"bouncin"}});
	s.conn.Send(&irc.Message{Command: "USER", Params: []string{"bouncin", "0", "*", "Bouncin test"}});
	return s;
}

func (server *server) recvFunc(msg *irc.Message) {
	events.Notify("RecvFromServer", server, msg)
}

func (server *server) errorFunc(err os.Error) {
	// TODO handle error
	log.Stderrf("server connection to %s failed: %s\n", server.conn.RemoteAddr(), err)
}

func (server *server) RemoteAddr() net.Addr {
	return server.conn.RemoteAddr()
}

func (server *server) Network() *Network {
	return server.network
}

// Send transmits a message by notifying SendToServer event chain.
func (server *server) Send(msg *irc.Message) {
	events.Notify("SendToServer", server, msg)
}

// sendToServer is the last handler in the SendToServer chain.  It performs the
// actual irc.Conn.Send() call which causes the message to be transmitted.
func sendToServer(conn Conn, msg *irc.Message) events.EventAction {
	conn.(*server).conn.Send(msg);
	return events.EventStop;
}

func init() {
	events.AddChain("RecvFromServer", InvokeSendRecv);
	events.AddChain("SendToServer", InvokeSendRecv);

	events.AddHandler("SendToServer", "server", events.PrioLast, sendToServer);
}
