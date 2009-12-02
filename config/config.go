// The config package provides configuration file handling.
package config

import (
	"io";
	"json";
	"log";
)

const CONFIG_FILENAME = ".bouncinrc"

// Plugins can keep persistent string settings.
type PluginData map[string] string;

// Top-level configuration.
type Config struct {
	Networks []Network;
	Extra PluginData;
}

// Per-network configuration.
type Network struct {
	Name string;
	Listen string;
	Server string;
	Channels []Channel;
	Extra PluginData;
}

// Per-channel configuration.
type Channel struct {
	Name string;
	Extra PluginData;
}

func ParseConfig(config string) *Config {
	var c = &Config{};
	ok, errtok := json.Unmarshal(config, c);
	if !ok {
		log.Exitf("Config syntax error: %s\n", errtok);
	}
	return c;
}

func ReadConfig() {
	content, err := io.ReadFile(CONFIG_FILENAME);
	if err != nil {
		log.Exitf("Can't open config file: %s\n", err);
	}
	ParseConfig(string(content));
}
