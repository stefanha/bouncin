package core

import (
	"net";
	"reflect";
	"irc";
	"events";
)

// A client or server IRC connection.
type Conn interface {
	// Send transmits an IRC message.
	Send(*irc.Message);

	// Network returns the network this connection belongs to.
	Network() *Network;

	// RemoteAddr returns the network address of the other end of the connection.
	RemoteAddr() net.Addr;
}

// InvokeSendRecv calls RecvFromClient, RecvFromServer, SendToClient, and
// SendToServer event handlers.  This function is public is that server and
// client modules can share it rather than having their own private copies.
func InvokeSendRecv(fn interface{}, args ...) events.EventAction {
	sendRecv	:= fn.(func(Conn, *irc.Message) events.EventAction);
	structValue	:= reflect.NewValue(args).(*reflect.StructValue);
	conn		:= structValue.Field(0).Interface().(Conn);
	msg		:= structValue.Field(1).Interface().(*irc.Message);
	return sendRecv(conn, msg);
}
