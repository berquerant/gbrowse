package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/berquerant/gbrowse/browse"
	"github.com/berquerant/gbrowse/ctxlog"
	"github.com/berquerant/gbrowse/git"
	"github.com/berquerant/gbrowse/parse"
	"github.com/berquerant/gbrowse/urlx"
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
)

type envConfig struct {
	Git     string `env:"GBROWSE_GIT" envDefault:"git"`
	IsDebug bool   `env:"GBROWSE_DEBUG"`
}

func (c *envConfig) Logger() (*zap.Logger, error) {
	if c.IsDebug {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
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

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		printOnly     = flag.Bool("print", false, "only print generated url")
		defaultBranch = flag.Bool("default", false, "use default branch instead of the current branch")
		config        envConfig
	)

	log.SetFlags(0)
	flag.Usage = Usage
	flag.Parse()

	var envList []map[string]any
	fail(env.ParseWithOptions(&config, env.Options{
		OnSet: func(tag string, value any, isDefault bool) {
			envList = append(envList, map[string]any{
				"tag":     tag,
				"value":   value,
				"default": isDefault,
			})
		},
	}))
	rawLogger, err := config.Logger()
	fail(err)
	for _, d := range envList {
		rawLogger.Debug("env",
			zap.Any("tag", d["tag"]),
			zap.Any("value", d["value"]),
			zap.Any("default", d["default"]),
		)
	}
	flag.VisitAll(func(f *flag.Flag) {
		rawLogger.Debug("flag",
			zap.String("tag", f.Name),
			zap.String("value", f.Value.String()),
			zap.Bool("default", f.Value.String() == f.DefValue),
		)
	})

	os.Exit(run(ctxlog.With(context.Background(), rawLogger), &args{
		config:        &config,
		target:        flag.Arg(0),
		printOnly:     *printOnly,
		defaultBranch: *defaultBranch,
	}))
}

type args struct {
	config        *envConfig
	target        string
	printOnly     bool
	defaultBranch bool
}

func run(ctx context.Context, args *args) int {
	logger := ctxlog.From(ctx)
	defer func() {
		_ = logger.Sync()
	}()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	target, err := parse.ReadTarget(args.target)
	if err != nil {
		logger.Error("parse target",
			zap.Error(err),
		)
		return 1
	}

	targetURL, err := urlx.Build(
		ctx,
		git.New(git.WithGitCommand(args.config.Git)),
		target,
		urlx.WithDefaultBranch(args.defaultBranch),
	)
	if err != nil {
		logger.Error("build url",
			zap.Error(err),
		)
		return 1
	}

	if args.printOnly {
		fmt.Print(targetURL)
		return 0
	}

	if err := browse.Run(ctx, targetURL); err != nil {
		logger.Error("browse",
			zap.Error(err),
		)
		return 1
	}
	return 0
}
