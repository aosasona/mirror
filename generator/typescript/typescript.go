package typescript

import "go.trulyao.dev/mirror/config"

type Config struct{}

func (c *Config) Name() string {
	return "typescript"
}

func (c *Config) Extension() string {
	return "typescript"
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

var _ config.LanguageConfigInterface = &Config{}
