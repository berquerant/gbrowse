package urlx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/berquerant/gbrowse/git"
	"github.com/berquerant/gbrowse/parse"
)

// Build assembles url from repository and specified path.
func Build(ctx context.Context, gitCommand git.Git, target *parse.Target) (string, error) {
	u, err := build(ctx, gitCommand, target)
	if err != nil {
		return "", fmt.Errorf("failed to build url: %w", err)
	}
	return u, nil
}

func build(ctx context.Context, gitCommand git.Git, target *parse.Target) (string, error) {
	type result struct {
		repoUrl  string
		branch   string
		path     string
		fragment string
	}
	var res result

	if err := func() error {
		{
			r, err := gitCommand.RemoteOriginUrl(ctx)
			if err != nil {
				return err
			}
			res.repoUrl = parse.ReadRepoUrl(r)
		}
		{
			r, err := gitCommand.HeadObjectName(ctx)
			if err != nil {
				return err
			}
			res.branch = r
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
		res.repoUrl, res.branch, res.path, res.fragment,
	), nil
}

func isDirectory(path string) (bool, error) {
	x, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return x.Mode().IsDir(), nil
}
