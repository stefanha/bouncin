package irc

type Message struct {
	Prefix	string;		// optional
	Command	string;
	Params	[]string;
}

func (msg *Message) String() string {
	s := "";

	if msg.Prefix != "" {
		s = ":" + msg.Prefix + " "
	}

	s += msg.Command;

	for i := 0; i < len(msg.Params); i++ {
		spacer := " ";
		if i == len(msg.Params) - 1 {
			spacer = " :"
		}
		s += spacer + msg.Params[i];
	}

	return s;
}

func Parse(line string) *Message {
	return &Message{};
}

func Fmt(command string, params ...) string {
	return "TODO"
}
