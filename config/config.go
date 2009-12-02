// The config package provides configuration file handling.
package config

import (
	"io";
	"json";
	"bytes";

	"log";
)

const (
	configFilename	= ".bouncinrc";
	configPerms	= 0600;
)

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
		log.Exitf("Config syntax error: %s\n", errtok)
	}
	return c;
}

func ReadConfig() *Config {
	content, err := io.ReadFile(configFilename);
	if err != nil {
		log.Exitf("Can't open config file: %s\n", err)
	}
	return ParseConfig(string(content));
}

func WriteConfig(config *Config) {
	buf := &bytes.Buffer{};
	err := json.Marshal(buf, *config);
	if err != nil {
		log.Exitf("Failed to marshal config file: %s\n", err)
	}

	// TODO atomic update without truncating file
	err = io.WriteFile(configFilename, buf.Bytes(), configPerms);
	if err != nil {
		log.Exitf("Can't write config file: %s\n", err)
	}
}
