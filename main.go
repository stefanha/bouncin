package main

import (
	"net";
	"runloop";
	"network";
)

func main() {
	addr, _ := net.ResolveTCPAddr("0.0.0.0:1234");
	listen, _ := net.ListenTCP("tcp", addr);
	network.Add("admin", nil, listen);

	runloop.Run();
}
