package irc

import "testing"

type tcase struct {
	line string;
	msg Message;
};
var cases = []tcase {
	tcase {
		":prefix command one two :three four",
		Message{"prefix", "command", []string{"one", "two", "three four"}},
	},
	tcase {
		"command",
		Message{Command: "command"},
	},
	tcase {
		"command :one",
		Message{Command: "command", Params: []string{"one"}},
	},
	tcase {
		"command one :two",
		Message{Command: "command", Params: []string{"one", "two"}},
	},
	tcase {
		"com:mand one :two three",
		Message{Command: "com:mand", Params: []string{"one", "two three"}},
	},
};

func TestMessageString(t *testing.T) {
	for i, tc := range cases {
		s := tc.msg.String();
		if s != tc.line {
			t.Errorf("Case %d: expected \"%s\", got \"%s\"", i, tc.line, s)
		}
	}
}

func TestParse(t *testing.T) {
	compare := func(a, b *Message) bool {
		if !(a.Prefix == b.Prefix && a.Command == b.Command && len(a.Params) == len(b.Params)) {
			return false;
		}
		for i, c := range a.Params {
			if c != b.Params[i] {
				return false;
			}
		}
		return true;
	};

	for i, tc := range cases {
		msg := Parse(tc.line);
		if msg == nil || !compare(&tc.msg, msg) {
			t.Errorf("Case %d: expected \"%s\", got \"%s\"", i, &tc.msg, msg);
		}
	}
}
