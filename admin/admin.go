// The admin package provides management commands.
package admin

import (
	"log";
	"core";
	"runloop";
	"plugins";
	"commands";
)

type admin struct {
}

func shutdown(conn core.Conn, argv []string) {
	log.Stderrf("shutdown command from %s\n", conn.RemoteAddr());
	runloop.Quit();
}

func (*admin) Enable() bool {
	commands.AddCommand("shutdown", shutdown,
`usage: shutdown
Terminates the program, disconnecting clients and leaving IRC networks.`);
	return true;
}

func (*admin) Disable() {
	commands.RemoveCommand("shutdown")
}

func init() {
	plugins.Register("admin", &admin{})
}
