package irc

import "strings"
import "regexp"
import "container/vector"

type Message struct {
	Prefix	string;		// optional
	Command string;
	Params	[]string;
}

func (msg *Message) String() string {
	s := "";

	if msg.Prefix != "" {
		s = ":" + msg.Prefix + " "
	}

	s += msg.Command;

	for i := 0; i < len(msg.Params); i++ {
		param	:= msg.Params[i];
		spacer	:= " ";
		if i == len(msg.Params) - 1 && strings.Index(param, " ") != -1 {
			spacer = " :"
		}
		s += spacer + param;
	}

	return s;
}

var re, _ = regexp.Compile("^(:[^ ]+ )?([^ ]+)( ?.*)$");

func Parse(line string) *Message {
	m := &Message{};
	matches := re.MatchStrings(line);
	if len(matches) != 4 {
		return nil;
	}

	prefix, command, params := matches[1], matches[2], matches[3];
	if len(prefix) > 2 {
		m.Prefix = prefix[1:len(prefix) - 1];
	}
	m.Command = strings.ToUpper(command);

	paramArray := vector.NewStringVector(0);
	param := "";
	for i, c := range params {
		switch {
		case c == ' ':
			if param != "" {
				paramArray.Push(param);
				param = "";
			}
		case c == ':' && param == "":
			paramArray.Push(params[i + 1:len(params)]);
			goto done;
		default:
			param += string(c);
		}
	}
	if param != "" {
		paramArray.Push(param);
		param = "";
	}
done:
	m.Params = paramArray.Data();
	return m;
}
