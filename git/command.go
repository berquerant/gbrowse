package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/berquerant/gbrowse/execx"
)

//go:generate go tool goconfig -field "GitCommand string" -option -output config_generated.go

// Git is git runner.
type Git interface {
	DefaultBranch(ctx context.Context) (string, error)
	RemoteOriginURL(ctx context.Context) (string, error)
	HeadObjectName(ctx context.Context) (string, error)
	ShowPrefix(ctx context.Context) (string, error)
	RelativePath(ctx context.Context, path string) (string, error)
	DescribeTag(ctx context.Context) (string, error)
	ShowCurrent(ctx context.Context) (string, error)
	CommitHash(ctx context.Context) (string, error)
}

type gitImpl struct {
	config *Config
}

func New(opt ...ConfigOption) Git {
	config := NewConfigBuilder().
		GitCommand("git").
		Build()
	for _, f := range opt {
		f(config)
	}
	return &gitImpl{
		config: config,
	}
}

func (g *gitImpl) DefaultBranch(ctx context.Context) (string, error) {
	r, err := g.run(ctx, "remote", "show", "origin")
	if err != nil {
		return "", err
	}

	for _, p := range strings.Split(r, "\n") {
		if strings.Contains(p, "HEAD branch:") {
			xs := strings.Split(p, ":")
			if len(xs) < 2 {
				break
			}
			if branch := strings.ReplaceAll(xs[1], " ", ""); branch != "" {
				return branch, nil
			}
		}
	}

	return "", fmt.Errorf("cannot find default branch from %s", r)
}

func (g *gitImpl) RemoteOriginURL(ctx context.Context) (string, error) {
	return g.run(ctx, "config", "--get", "remote.origin.url")
}

func (g *gitImpl) HeadObjectName(ctx context.Context) (string, error) {
	return g.run(ctx, "rev-parse", "--abbrev-ref", "@")
}

func (g *gitImpl) ShowPrefix(ctx context.Context) (string, error) {
	return g.run(ctx, "rev-parse", "--show-prefix")
}

func (g *gitImpl) RelativePath(ctx context.Context, path string) (string, error) {
	return g.run(ctx, "ls-files", "--full-name", path)
}

func (g *gitImpl) DescribeTag(ctx context.Context) (string, error) {
	return g.run(ctx, "describe", "--tags", "--abbrev=0")
}

func (g *gitImpl) ShowCurrent(ctx context.Context) (string, error) {
	return g.run(ctx, "branch", "--show-current")
}

func (g *gitImpl) CommitHash(ctx context.Context) (string, error) {
	return g.run(ctx, "rev-parse", "@")
}

func (g *gitImpl) run(ctx context.Context, arg ...string) (string, error) {
	return execx.Run(ctx, g.config.GitCommand.Get(), arg...)
}
