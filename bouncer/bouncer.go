// The bouncer package implements an IRC proxy.
package bouncer

import (
	"irc";
	"plugins";
	"events";
	"core";
)

type bouncer struct {
}

func recvFromClient(client core.Conn, msg *irc.Message) events.EventAction {
	// TODO any filtering/mangling necessary here?
	// TODO need to send to other clients?
	client.Network().SendToServer(msg);
	return events.EventStop;
};

func recvFromServer(server core.Conn, msg *irc.Message) events.EventAction {
	server.Network().SendToClients(msg);
	return events.EventStop;
};

func (bouncer *bouncer) Enable() bool {
	events.AddHandler("RecvFromClient", "bouncer", events.PrioLast, recvFromClient);
	events.AddHandler("RecvFromServer", "bouncer", events.PrioLast, recvFromServer);
	return true;
}

func (bouncer *bouncer) Disable() {
	events.RemoveHandler("RecvFromClient", "bouncer");
	events.RemoveHandler("RecvFromServer", "bouncer");
}

func init() {
	plugins.Register("bouncer", &bouncer{})
}
