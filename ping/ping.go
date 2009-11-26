// The ping package responds to IRC ping commands.
package ping

// TODO send ping commands after a connection goes quiet

import (
	"irc";
	"plugins";
	"events";
	"core";
)

type ping struct {
}

func handlePing(conn core.Conn, pingMsg *irc.Message) events.EventAction {
	if pingMsg.Command != "PING" {
		return events.EventContinue;
	}

	conn.Send(&irc.Message{Command: "PONG", Params: []string{"bouncin"}}); // TODO use proper server name param
	return events.EventStop;
}

func (ping *ping) Enable() bool {
	events.AddHandler("RecvFromClient", "ping", events.PrioNormal, handlePing);
	events.AddHandler("RecvFromServer", "ping", events.PrioNormal, handlePing);
	return true;
}

func (ping *ping) Disable() {
	events.RemoveHandler("RecvFromClient", "ping");
	events.RemoveHandler("RecvFromServer", "ping");
}

func init() {
	plugins.Register("ping", &ping{})
}
