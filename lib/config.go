package lib

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Command           string   `yaml:"command"`
	ProtectedContexts []string `yaml:"protectedContexts"`
	ConfirmString     string   `yaml:"confirmString"`
}

func parseConfig(conf []byte) (Config, error) {
	cfg := Config{}
	err := yaml.Unmarshal(conf, &cfg)
	return cfg, err
}
