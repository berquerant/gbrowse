package git

import (
	"context"

	"github.com/berquerant/gbrowse/execx"
)

//go:generate go run github.com/berquerant/goconfig@v0.3.0 -field "GitCommand string" -option -output config_generated.go

// Git is git runner.
type Git interface {
	RemoteOriginUrl(ctx context.Context) (string, error)
	HeadObjectName(ctx context.Context) (string, error)
	ShowPrefix(ctx context.Context) (string, error)
	RelativePath(ctx context.Context, path string) (string, error)
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

func (g *gitImpl) RemoteOriginUrl(ctx context.Context) (string, error) {
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

func (g *gitImpl) run(ctx context.Context, arg ...string) (string, error) {
	return execx.Run(ctx, g.config.GitCommand.Get(), arg...)
}
