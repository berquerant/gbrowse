package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/berquerant/gbrowse/browse"
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

  gbrowse -commit [commit] opens the commit page.
  If commit is omitted, the current commit is opened.

  gbrowse -compare <COMPARE> [target] opens the comparison page between COMPARE and target.
  If target is omitted, the current commit is used.

Environment variables:
  GIT
    git command, default is git.

  DEBUG
    enable debug log if set.

Flags:`

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	var (
		printOnly = flag.Bool("print", false, "only print generated url")
		commit    = flag.Bool("commit", false, "open commit page")
		compare   = flag.String("compare", "", "open compare page between the specified ref and target")
		envConfig = newEnvConfig()
		logger    = envConfig.logger()
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

	run(ctxlog.With(context.Background(), logger), &args{
		envConfig: envConfig,
		target:    flag.Arg(0),
		printOnly: *printOnly,
		commit:    *commit,
		compare:   *compare,
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
	envConfig *envConfig
	target    string
	printOnly bool
	commit    bool
	compare   string
}

func run(ctx context.Context, args *args) exitCode {
	logger := ctxlog.From(ctx)

	if args.commit && args.compare != "" {
		logger.Error("commit and compare cannot be used together")
		return eFailure
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	gitCommand := git.New(git.WithGitCommand(args.envConfig.Git))

	var (
		targetURL string
		err       error
	)

	switch {
	case args.commit:
		targetURL, err = urlx.BuildCommit(ctx, gitCommand, args.target)
	case args.compare != "":
		targetURL, err = urlx.BuildCompare(ctx, gitCommand, args.compare, args.target)
	default:
		var target *parse.Target
		target, err = parse.ReadTarget(args.target)
		if err != nil {
			logger.Error("parse target",
				ctxlog.Err(err),
			)
			return eFailure
		}
		targetURL, err = urlx.Build(
			ctx,
			gitCommand,
			target,
		)
	}

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
