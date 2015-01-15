package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	confValues map[string]toml.Primitive
	md         toml.MetaData
}

func NewConfig(filename string) (*Config, error) {
	c := &Config{
		make(map[string]toml.Primitive),
		toml.MetaData{},
	}

	var err error
	c.md, err = toml.DecodeFile(filename, c.confValues)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Load(section string, conf interface{}) error {
	if v, ok := c.confValues[section]; ok {
		return c.md.PrimitiveDecode(v, conf)
	}
	return fmt.Errorf("Config section %q missing", section)
}
