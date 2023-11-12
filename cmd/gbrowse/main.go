package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/berquerant/gbrowse/browse"
	"github.com/berquerant/gbrowse/config"
	"github.com/berquerant/gbrowse/ctxlog"
	"github.com/berquerant/gbrowse/env"
	"github.com/berquerant/gbrowse/git"
	"github.com/berquerant/gbrowse/parse"
	"github.com/berquerant/gbrowse/urlx"
	"golang.org/x/exp/slog"
)

type envConfig struct {
	Git     string
	IsDebug bool
}

func newEnvConfig() *envConfig {
	var c envConfig
	c.Git = env.GetOr("GBROWSE_GIT", "git")
	c.IsDebug = env.GetOr("GBROWSE_DEBUG", "") != ""
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

func parseConfig(filePath string) (*config.Config, error) {
	if filePath == "" {
		return config.Default(), nil
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return config.Parse(b)
}

const usage = `gbrowse - Open the repo in the browser

Usage:
  gbrowse [flags] [target]

  The target is PATH or FILE:LINUM.
  gbrowse PATH opens the PATH of the repo.
  gbrowse FILE:LINUM opens the line LINUM of the FILE of the repo.
  gbrowse opens the directory of the repo.

Environment variables:
  GBROWSE_GIT
    git command, default is git.

  GBROWSE_DEBUG
    enable debug log if set.

Flags:`

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	var (
		printOnly  = flag.Bool("print", false, "only print generated url")
		configFile = flag.String("config", "", "config file")
		envConfig  = newEnvConfig()
		logger     = envConfig.logger()
	)

	flag.Usage = Usage
	flag.Parse()

	logger.Debug("env", ctxlog.Any("values", envConfig))
	flag.VisitAll(func(f *flag.Flag) {
		logger.Debug("flag",
			ctxlog.S("tag", f.Name),
			ctxlog.S("value", f.Value.String()),
			ctxlog.B("default", f.Value.String() == f.DefValue),
		)
	})

	config, err := parseConfig(*configFile)
	if err != nil {
		logger.Error("parse config", ctxlog.Err(err))
		os.Exit(int(eFailure))
	}

	run(ctxlog.With(context.Background(), logger), &args{
		config:    config,
		envConfig: envConfig,
		target:    flag.Arg(0),
		printOnly: *printOnly,
	}).exit()
}

type exitCode int

const (
	eSuccess exitCode = iota
	eFailure
)

func (c exitCode) exit() {
	os.Exit(int(c))
}

type args struct {
	config    *config.Config
	envConfig *envConfig
	target    string
	printOnly bool
}

func run(ctx context.Context, args *args) exitCode {
	logger := ctxlog.From(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	target, err := parse.ReadTarget(args.target)
	if err != nil {
		logger.Error("parse target",
			ctxlog.Err(err),
		)
		return eFailure
	}

	gitCommand := git.New(git.WithGitCommand(args.envConfig.Git))
	phaseExecutor := urlx.NewPhaseExecutor(gitCommand)

	targetURL, err := urlx.Build(
		ctx,
		gitCommand,
		target,
		phaseExecutor,
		urlx.WithPhases(args.config.Phases),
	)
	if err != nil {
		logger.Error("build url",
			ctxlog.Err(err),
		)
		return eFailure
	}

	if args.printOnly {
		fmt.Print(targetURL)
		return eSuccess
	}

	if err := browse.Run(ctx, targetURL); err != nil {
		logger.Error("browse",
			ctxlog.Err(err),
		)
		return eFailure
	}
	return eSuccess
}
