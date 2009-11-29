package irc

import "fmt"

// A numeric IRC response message.
type Numeric struct {
	num	string;
	txt	string;
}

var RplWelcome	= &Numeric{"001", "Welcome to the Internet Relay Network %s!%s@%s"}
var RplYourHost	= &Numeric{"002", "Your host is %s, running version %s"}
var RplCreated	= &Numeric{"003", "This server was created %s"}
var RplMyInfo  = &Numeric{"004", "%s %s %s %s"}

func (numeric *Numeric) Message(nick string, args ...) *Message {
	txt := fmt.Sprintf(numeric.txt, args);
	return &Message{Command: numeric.num, Params: []string{nick, txt}};
}
