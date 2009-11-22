// The client package implements an IRC connection for a client.
package client

import (
	"os";
	"net";
	"log";
	"irc";
)

type Client struct {
	conn		*irc.Conn;

	// Login nick and real name
	nick		string;
	realname	string;
}

func New(conn net.Conn) *Client {
	var client *Client;

	recvFunc := func(msg *irc.Message) {
		client.recvFunc(msg)
	};
	errorFunc := func(err os.Error) {
		client.errorFunc(err)
	};

	client = &Client{conn: irc.NewConn(conn, recvFunc, errorFunc)};
	return client;
}

func (client *Client) recvFunc(msg *irc.Message) {
	// Connection registration is a special state, don't process messages
	// until the client gives its nick and user.
	if client.nick == "" {
		if msg.Command == "NICK" && len(msg.Params) == 1 {
			client.nick = msg.Params[0]
		}
		return;
	}
	if client.realname == "" {
		if msg.Command == "USER" && len(msg.Params) == 4 {
			client.realname = msg.Params[3];
			client.welcome();
		}
		return;
	}

	if msg.Command == "PING" {
		client.conn.Send(&irc.Message{Command: "PONG", Params: []string{"bouncin"}})
	}
}

func (client *Client) errorFunc(err os.Error) {
	// TODO handle error
	log.Stderrf("client %s failed: %s\n", client.conn.RemoteAddr(), err);
}

func (client *Client) welcome() {
	client.conn.Send(irc.RplWelcome.Message(client.nick, client.nick, client.nick, client.conn.RemoteAddr()));
	client.conn.Send(irc.RplYourHost.Message(client.nick, "bouncin", "0.1"));
	client.conn.Send(irc.RplCreated.Message(client.nick, "..."));
	client.conn.Send(irc.RplMyInfo.Message(client.nick, "bouncin", "0.1", "", ""));
}
