package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	kFile = "kafka.yml"
)

type Config struct {
	Brokers string
	Topic   string
}

func ParseConfig() (Config, error) {
	c := Config{}

	yml, err := ioutil.ReadFile(kFile)
	if err != nil {
		return c, err
	}

	return c, yaml.Unmarshal(yml, &c)
}
