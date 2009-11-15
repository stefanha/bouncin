package config

import "json"
import "log"
import "io"

const CONFIG_FILENAME = ".bouncinrc"

type Config struct {
	Networks []Network;
}

type Network struct {
	Name string;
	Listen string;
	Server string;
}

func ReadConfig() {
	content, err := io.ReadFile(CONFIG_FILENAME);
	if err != nil {
		log.Exitf("Can't open config file: %s\n", err);
	}
	ParseConfig(string(content));
}

func ParseConfig(config string) *Config {
	var c = Config{};
	ok, errtok := json.Unmarshal(config, &c);
	if !ok {
		log.Exitf("Config syntax error: %s\n", errtok);
	}
	return &c;
}
