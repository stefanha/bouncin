package core

import (
	"os";
	"net";
	"log";
	"irc";
	"runloop";
	"events";
)

// A clients implements a connection to an IRC client.
type client struct {
	conn		*irc.Conn;
	network		*Network;

	// Login nick and real name
	nick		string;
	realname	string;
}

// newclient returns a new client for a given network connection.
func newclient(conn net.Conn, network *Network) *client {
	var c *client;

	recvFunc := func(msg *irc.Message) {
		runloop.CallLater(func() { c.recvFunc(msg) })
	};
	errorFunc := func(err os.Error) {
		runloop.CallLater(func() { c.errorFunc(err) })
	};

	c = &client{conn: irc.NewConn(conn, recvFunc, errorFunc), network: network};
	return c;
}

func (c *client) errorFunc(err os.Error) {
	// TODO handle error
	log.Stderrf("client %s failed: %s\n", c.conn.RemoteAddr(), err);
}

func (c *client) welcome() {
	c.conn.Send(irc.RplWelcome.Message(c.nick, c.nick, c.nick, c.conn.RemoteAddr()));
	c.conn.Send(irc.RplYourHost.Message(c.nick, "bouncin", "0.1"));
	c.conn.Send(irc.RplCreated.Message(c.nick, "..."));
	c.conn.Send(irc.RplMyInfo.Message(c.nick, "bouncin", "0.1", "", ""));
}

func (c *client) recvFunc(msg *irc.Message) {
	// Connection registration is a special state, don't process messages
	// until the client gives its nick and user.
	if c.nick == "" {
		if msg.Command == "NICK" && len(msg.Params) == 1 {
			c.nick = msg.Params[0]
		}
		return;
	}
	if c.realname == "" {
		if msg.Command == "USER" && len(msg.Params) == 4 {
			c.realname = msg.Params[3];
			c.welcome();
		}
		return;
	}

	// The normal code path notifies the RecvFromclient event chain.
	events.Notify("RecvFromclient", c, msg);
}

func (c *client) Network() *Network {
	return c.network
}

func (c *client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *client) Send(msg *irc.Message) {
	events.Notify("SendToclient", c, msg)
}

func sendToclient(conn Conn, msg *irc.Message) events.EventAction {
	conn.(*client).conn.Send(msg);
	return events.EventStop;
}

func init() {
	events.AddChain("RecvFromclient", InvokeSendRecv);
	events.AddChain("SendToclient", InvokeSendRecv);

	events.AddHandler("SendToclient", "client", events.PrioLast, sendToclient);
}
