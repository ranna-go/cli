package config

import (
	"os"
	"path"

	"github.com/ranna-go/ranna/pkg/client"
	"github.com/traefik/paerser/env"
	"github.com/traefik/paerser/file"
)

type Config struct {
	client.Options
}

func Parse() (c Config, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	c.Endpoint = "https://public.ranna.zekro.de"
	c.Version = "v1"

	err = parseFiles(&c,
		path.Join(home, ".ranna", "config.json"),
		path.Join(home, ".ranna", "config.yaml"),
		path.Join(home, ".ranna", "config.yml"),
		path.Join(home, ".ranna", "config.toml"),
		"config.json",
		"config.yaml",
		"config.yml",
		"config.toml",
	)
	if err != nil {
		return
	}

	if err = env.Decode(os.Environ(), "RANNA_", &c); err != nil {
		return
	}

	return
}

func parseFile(path string, cfg *Config) (err error) {
	err = file.Decode(path, cfg)
	if os.IsNotExist(err) {
		err = nil
	}
	return
}

func parseFiles(cfg *Config, path ...string) (err error) {
	for _, p := range path {
		if err = parseFile(p, cfg); err != nil {
			return
		}
	}
	return
}
