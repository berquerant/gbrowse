package urlx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/berquerant/gbrowse/config"
	"github.com/berquerant/gbrowse/git"
	"github.com/berquerant/gbrowse/parse"
)

//go:generate go run github.com/berquerant/goconfig@v0.3.0 -field "Phases []config.Phase" -option -output config_generated.go

// Build assembles url from repository and specified path.
func Build(ctx context.Context, gitCommand git.Git, target *parse.Target, executor PhaseExecutor, opt ...ConfigOption) (string, error) {
	config := NewConfigBuilder().
		Phases([]config.Phase{}).
		Build()
	config.Apply(opt...)

	doPhases := func(ctx context.Context) (string, error) {
		return ExecutePhases(ctx, config.Phases.Get(), executor)
	}
	u, err := build(ctx, gitCommand, target, doPhases)
	if err != nil {
		return "", fmt.Errorf("failed to build url: %w", err)
	}
	return u, nil
}

func build(ctx context.Context, gitCommand git.Git, target *parse.Target, doPhases func(context.Context) (string, error)) (string, error) {
	type result struct {
		repoURL  string
		ref      string
		path     string
		fragment string
	}
	var res result

	if err := func() error {
		{
			r, err := gitCommand.RemoteOriginURL(ctx)
			if err != nil {
				return err
			}
			res.repoURL = parse.ReadRepoURL(r)
		}
		{
			r, err := doPhases(ctx)
			if err != nil {
				return err
			}
			res.ref = r
		}
		{
			if isDir, err := isDirectory(target.Path()); err != nil || isDir {
				r, err := gitCommand.ShowPrefix(ctx)
				if err != nil {
					return err
				}
				res.path = filepath.Join(r, target.Path())
			} else {
				r, err := gitCommand.RelativePath(ctx, target.Path())
				if err != nil {
					return err
				}
				res.path = r
			}
		}
		{
			if linum, ok := target.Linum(); ok {
				res.fragment = fmt.Sprintf("#L%d", linum)
			}
		}
		return nil
	}(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/blob/%s/%s%s",
		res.repoURL, res.ref, res.path, res.fragment,
	), nil
}

func isDirectory(path string) (bool, error) {
	x, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return x.Mode().IsDir(), nil
}
