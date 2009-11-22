package server

import (
	"os";
	"net";
	"log";
	"irc";
)

// A Server implements a connection to an IRC server.
type Server struct {
	conn		*irc.Conn;
}

// New returns a new Server for a given network connection.
func New(conn net.Conn) *Server {
	var server *Server;

	if conn == nil {
		return nil
	}

	recvFunc := func(msg *irc.Message) {
		server.recvFunc(msg)
	};
	errorFunc := func(err os.Error) {
		server.errorFunc(err)
	};

	server = &Server{irc.NewConn(conn, recvFunc, errorFunc)};
	server.conn.Send(&irc.Message{Command: "NICK", Params: []string{"bouncin"}});
	server.conn.Send(&irc.Message{Command: "USER", Params: []string{"bouncin", "0", "*", "Bouncin test"}});
	return server;
}

func (server *Server) recvFunc(msg *irc.Message) {
}

func (server *Server) errorFunc(err os.Error) {
	// TODO handle error
	log.Stderrf("server connection to %s failed: %s\n", server.conn.RemoteAddr(), err)
}
