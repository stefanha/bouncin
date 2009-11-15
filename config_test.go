package config

import "testing"

func compareNetwork(a, b *Network) bool {
	return a.Name == b.Name && a.Listen == b.Listen && a.Server == b.Server;
}

func compareConfig(a, b *Config) bool {
	if len(a.Networks) != len(b.Networks) {
		return false;
	}
	for i, network := range a.Networks {
		if !compareNetwork(&network, &b.Networks[i]) {
			return false;
		}
	}
	return true;
}

type tcase struct {
	configString string;
	config *Config;
}

var tcases = []tcase {
	tcase {
		`{"networks": [{"name": "freenode", "listen": "0.0.0.0:1234", "server": "chat.freenode.net"}]}`,
		&Config{Networks: []Network{
			Network{"freenode", "0.0.0.0:1234", "chat.freenode.net"},
		}}
	},
}

func TestConfig(t *testing.T) {
	for i, tc := range tcases {
		config := ParseConfig(tc.configString);
		if !compareConfig(config, tc.config) {
			t.Errorf("Case %d: expected %s, got %s\n", i, tc.config, config);
		}
	}
}
