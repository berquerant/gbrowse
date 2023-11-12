package main

import (
	"github.com/berquerant/gbrowse/config"
)

func parseConfig(envFileOrPath, configOrFile string) *config.Config {
	return config.ParseStringOrFile(configOrFile, config.ParseStringOrFile(envFileOrPath, nil))
}
