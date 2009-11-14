package irc

import "testing"

func TestMessageString(t *testing.T) {
	type tcase struct {
		msg Message;
		expected string;
	};
	cases := []tcase {
		tcase {
			Message{Command: "command", Params: []string{}},
			"command"
		},
		tcase {
			Message{Command: "command", Params: []string{"one"}},
			"command :one"
		},
		tcase {
			Message{Command: "command", Params: []string{"one", "two"}},
			"command one :two"
		},
		tcase {
			Message{"prefix", "command", []string{"one", "two", "three four"}},
			":prefix command one two :three four"
		},
	};

	for i, tc := range cases {
		s := tc.msg.String();
		if s != tc.expected {
			t.Errorf("Case %d: expected \"%s\", got \"%s\"", i, tc.expected, s)
		}
	}
}
