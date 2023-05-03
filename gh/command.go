package gh

import (
	"context"
	"fmt"

	"github.com/berquerant/gbrowse/execx"
)

//go:generate go run github.com/berquerant/goconfig@v0.3.0 -field "GHCommand string" -option -output config_generated.go

// GH is github cli runner.
type GH interface {
	Browse(ctx context.Context, target string) error
}

type ghImpl struct {
	config *Config
}

func New(opt ...ConfigOption) GH {
	config := NewConfigBuilder().
		GHCommand("gh").
		Build()
	for _, f := range opt {
		f(config)
	}
	return &ghImpl{
		config: config,
	}
}

func (g *ghImpl) Browse(ctx context.Context, target string) error {
	arg := []string{"browse"}
	if target != "" {
		arg = append(arg, target)
	}
	if _, err := execx.Run(ctx, g.config.GHCommand.Get(), arg...); err != nil {
		return fmt.Errorf("failed to github browse: %w", err)
	}
	return nil
}
