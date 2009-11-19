package main

import (
	"net";
	"runloop";
	"plugins";
	"network";

	// Plugins are imported so their init() functions are called
	_ "admin";
)

func main() {
	addr, _ := net.ResolveTCPAddr("0.0.0.0:1234");
	listen, _ := net.ListenTCP("tcp", addr);
	network.Add("admin", nil, listen);

	plugins.Enable("admin");

	runloop.Run();
}
