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
	var (
		repoURL  string
		ref      string
		path     string
		fragment string
	)
	if err := func() error {
		var err error
		if repoURL, err = gitCommand.RemoteOriginURL(ctx); err != nil {
			return err
		}
		if ref, err = gitCommand.CommitHash(ctx); err != nil {
			return err
		}
		if isDir, err := isDirectory(target.Path()); err != nil || isDir {
			r, err := gitCommand.ShowPrefix(ctx)
			if err != nil {
				return err
			}
			path = filepath.Join(r, target.Path())
		} else if path, err = gitCommand.RelativePath(ctx, target.Path()); err != nil {
			return err
		}
		if linum, ok := target.Linum(); ok {
			fragment = fmt.Sprintf("#L%d", linum)
		}
		return nil
	}(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/blob/%s/%s%s",
		repoURL, ref, path, fragment,
	), nil
}

func isDirectory(path string) (bool, error) {
	x, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return x.Mode().IsDir(), nil
}
