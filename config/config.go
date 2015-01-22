package config

import (
	"os"

	"../parser"
)

type Config struct {
	env *parser.Env
}

var DefaultConfigLocation = os.ExpandEnv("$HOME/.gig/config")

func NewConfig(filename string) (*Config, error) {
	c := &Config{}
	var err error
	c.env, err = parser.NewBootstrappedEnv()
	if err != nil {
		return nil, err
	}

	tune, err := c.env.LoadTune(filename)
	if err != nil {
		// Well that didn't work. Let's try something else.
		c, err = NewDefaultConfig()
		if err != nil {
			// Welp, we tried.
			return nil, err
		}
	}

	_, err = c.env.Eval(tune)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewDefaultConfig() (*Config, error) {
	c := &Config{}
	var err error
	c.env, err = parser.NewBootstrappedEnv()
	if err != nil {
		return nil, err
	}

	tune, err := c.env.LoadTune(DefaultConfigLocation)
	if err != nil {
		return nil, err
	}

	_, err = c.env.Eval(tune)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Get(key string) (interface{}, error) {
	return c.env.Get(key)
}
