package main

import (
	"net";
	"runloop";
	"plugins";
	"core";

	// Plugins are imported so their init() functions are called
	_ "admin";
	_ "bouncer";
	_ "ping";
)

func main() {
	addr, _		:= net.ResolveTCPAddr("0.0.0.0:1234");
	listen, _	:= net.ListenTCP("tcp", addr);
	server, _	:= net.Dial("tcp", "", "chat.freenode.net:6667");

	core.AddNetwork("freenode", server, listen);

	plugins.Enable("admin");
	plugins.Enable("bouncer");
	plugins.Enable("ping");

	runloop.Run();
}
