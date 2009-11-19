package main

import (
	"net";
	"runloop";
)

func main() {
	addr, _ := net.ResolveTCPAddr("0.0.0.0:1234");
	listen, _ := net.ListenTCP("tcp", addr);
	NewNetwork("admin", nil, listen);

	runloop.Run();
}
