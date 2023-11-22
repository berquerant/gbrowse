package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/berquerant/gbrowse/browse"
	"github.com/berquerant/gbrowse/config"
	"github.com/berquerant/gbrowse/ctxlog"
	"github.com/berquerant/gbrowse/git"
	"github.com/berquerant/gbrowse/parse"
	"github.com/berquerant/gbrowse/urlx"
)

const usage = `gbrowse - Open the repo in the browser

Usage:
  gbrowse [flags] [target]

  The target is PATH or FILE:LINUM.
  gbrowse PATH opens the PATH of the repo.
  gbrowse FILE:LINUM opens the line LINUM of the FILE of the repo.
  gbrowse opens the directory of the repo.

Config:

  {
    "phases": [
      PHASE, ...
    ],
    "defs": [
      {
        "id": ID,
        "cmd": ["command", ...]
      }, ...
    ]
  }

phases determines the search order for ref (commit, branch, tag).
PHASE is branch, default_branch, tag, commit or id in def.
defs is custom phases, cmd should return a string like commit hash, for example,

  {
    "phases": ["echo-master"],
    "defs": [{"id": "echo-master", "cmd": ["echo", "master"]}]
  }

sets ref to "master".
PHASE can also be specified by -phase flag.

If all searches fail, search commit.

Environment variables:
  GBROWSE_GIT
    git command, default is git.

  GBROWSE_DEBUG
    enable debug log if set.

  GBROWSE_CONFIG
    config file or string.
    -config overwrites this.

Flags:`

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	var (
		printOnly    = flag.Bool("print", false, "only print generated url")
		configOrFile = flag.String("config", "", "config or file")
		phases       []config.Phase
		envConfig    = newEnvConfig()
		logger       = envConfig.logger()
	)
	flag.Func("phase", "phases separated by comma", func(s string) error {
		ss := strings.Split(s, ",")
		phases = make([]config.Phase, len(ss))
		for i, x := range ss {
			phases[i] = config.NewPhase(x)
		}
		return nil
	})

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

	run(ctxlog.With(context.Background(), logger), &args{
		configOrFile: *configOrFile,
		envConfig:    envConfig,
		phases:       phases,
		target:       flag.Arg(0),
		printOnly:    *printOnly,
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
	configOrFile string
	envConfig    *envConfig
	phases       []config.Phase
	target       string
	printOnly    bool
}

func run(ctx context.Context, args *args) exitCode {
	logger := ctxlog.From(ctx)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	config := parseConfig(args.envConfig.Config, args.configOrFile)
	if len(args.phases) > 0 {
		config.Phases = args.phases
	}

	target, err := parse.ReadTarget(args.target)
	if err != nil {
		logger.Error("parse target",
			ctxlog.Err(err),
		)
		return eFailure
	}

	gitCommand := git.New(git.WithGitCommand(args.envConfig.Git))
	customExecutor := urlx.NewCustomPhaseExecutor(config.Definitions)
	phaseExecutor := urlx.NewPhaseExecutor(gitCommand, customExecutor)

	targetURL, err := urlx.Build(
		ctx,
		gitCommand,
		target,
		phaseExecutor,
		urlx.WithPhases(config.Phases),
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
