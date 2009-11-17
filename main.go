package main

import "net";

func main() {
	addr, _ := net.ResolveTCPAddr("0.0.0.0:1234");
	listen, _ := net.ListenTCP("tcp", addr);
	NewNetwork("admin", nil, listen);

	RunDispatcher();
}
