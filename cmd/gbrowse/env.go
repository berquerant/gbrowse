package main

import (
	"log/slog"
	"os"

	"github.com/berquerant/gbrowse/ctxlog"
	"github.com/berquerant/gbrowse/env"
)

type envConfig struct {
	Git     string
	IsDebug bool
	Config  string
}

func newEnvConfig() *envConfig {
	var c envConfig
	c.Git = env.GetOr("GBROWSE_GIT", "git")
	c.IsDebug = env.GetOr("GBROWSE_DEBUG", "") != ""
	c.Config = env.GetOr("GBROWSE_CONFIG", "")
	return &c
}

func (c *envConfig) logLevel() slog.Level {
	if c.IsDebug {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func (c *envConfig) logHandlerOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: c.logLevel(),
	}
}

func (c *envConfig) logger() ctxlog.Logger {
	return ctxlog.New(slog.New(slog.NewJSONHandler(os.Stdout, c.logHandlerOptions())))
}
