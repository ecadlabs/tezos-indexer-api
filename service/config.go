package service

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MaxConnections int           `yaml:"max_connections"`
	Timeout        time.Duration `yaml:"timeout"`
	URI            string        `yaml:"db_uri"`
	HTTPAddress    string        `yaml:"http_address"`
	LogHTTP        bool          `yaml:"log_http"`
}

func (c *Config) Load(name string) error {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, c)
}
