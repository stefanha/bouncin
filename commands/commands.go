// The commands package 
package commands

import (
	"fmt";
	"sort";
	"strings";
	"log";
	"irc";
	"events";
	"core";
)

// Commands receive the arguments array and the source connection.
type CommandFunc func(conn core.Conn, argv []string)

type command struct {
	name	string;
	fn	CommandFunc;
	help	string;
}

func sendLines(conn core.Conn, lines []string) {
	// TODO should we use INFO or something else?
	// TODO should the real nick be used?
	for _, line := range lines {
		conn.Send(irc.RplInfo.Message("bouncin", line))
	}
	conn.Send(irc.RplEndOfInfo.Message("bouncin"));
}

func (cmd *command) showHelp(conn core.Conn) {
	sendLines(conn, strings.Split(cmd.help, "\n", 0))
}

// Registered commands are stored as a map.
var cmds = make(map[string] *command)

// AddCommand registers a command.
func AddCommand(name string, fn CommandFunc, help string) {
	log.Stderrf("AddCommand %s\n", name);
	cmds[name] = &command{name, fn, help};
}

// RemoveCommand unregisters a command.
func RemoveCommand(name string) {
	log.Stderrf("RemoveCommand %s\n", name);
	cmds[name] = nil, false;
}

func listCommands(conn core.Conn) {
	// Sort command names alphabetically
	names := make([]string, len(cmds));
	i := 0;
	for name := range cmds {
		names[i] = name;
		i++;
	}
	sort.SortStrings(names);

	// Print command names
	lines := make([]string, (len(names) + 3) / 4);
	for i = 0; i < len(names); i += 4 {
		a, b, c, d := names[i], "", "", "";
		if i + 1 < len(names) {
			b = names[i + 1]
		}
		if i + 2 < len(names) {
			c = names[i + 2]
		}
		if i + 3 < len(names) {
			d = names[i + 3]
		}
		lines[i / 4] = fmt.Sprintf("%20s%20s%20s%20s", a, b, c, d);
	}
	sendLines(conn, lines);
}

// The help command.
func help(conn core.Conn, argv []string) {
	var cmd *command;
	var ok bool;
	if len(argv) != 2 {
		listCommands(conn)
	} else if cmd, ok = cmds[argv[1]]; !ok {
		listCommands(conn)
	} else {
		cmd.showHelp(conn)
	}
}

// recvFromClient dispatches commands received from clients.
func recvFromClient(conn core.Conn, msg *irc.Message) events.EventAction {
	if msg.Command != "BOUNCIN" {
		return events.EventContinue
	}

	if len(msg.Params) == 0 {
		return events.EventStop
	}

	if cmd, ok := cmds[msg.Params[0]]; ok {
		cmd.fn(conn, msg.Params)
	}
	return events.EventStop;
}

func init() {
	events.AddHandler("RecvFromClient", "commands", events.PrioNormal, recvFromClient);

	AddCommand("help", help,
`usage: help [<command>]
Shows documentation on a given command.  If no command is given, then all
available commands are listed.`);
}
