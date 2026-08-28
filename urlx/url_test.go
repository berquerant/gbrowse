package urlx_test

import (
	"context"
	"testing"

	"github.com/berquerant/gbrowse/parse"
	"github.com/berquerant/gbrowse/urlx"
	"github.com/stretchr/testify/assert"
)

type mockGit struct {
	remoteOriginURL string
	commitHash      string
	showPrefix      string
	relativePath    string
}

func (m *mockGit) DefaultBranch(ctx context.Context) (string, error)   { return "main", nil }
func (m *mockGit) RemoteOriginURL(ctx context.Context) (string, error) { return m.remoteOriginURL, nil }
func (m *mockGit) HeadObjectName(ctx context.Context) (string, error)  { return "HEAD", nil }
func (m *mockGit) ShowPrefix(ctx context.Context) (string, error)      { return m.showPrefix, nil }
func (m *mockGit) RelativePath(ctx context.Context, path string) (string, error) {
	return m.relativePath, nil
}
func (m *mockGit) DescribeTag(ctx context.Context) (string, error) { return "v0.1.0", nil }
func (m *mockGit) ShowCurrent(ctx context.Context) (string, error) { return "main", nil }
func (m *mockGit) CommitHash(ctx context.Context) (string, error)  { return m.commitHash, nil }

func TestBuild(t *testing.T) {
	gitCmd := &mockGit{
		remoteOriginURL: "git@github.com:berquerant/gbrowse.git",
		commitHash:      "0123456789abcdef",
		showPrefix:      "sub/dir/",
		relativePath:    "sub/dir/main.go",
	}

	for _, tc := range []struct {
		name    string
		target  string
		want    string
		wantErr bool
	}{
		{
			name:   "file with line number",
			target: "main.go:42",
			want:   "https://github.com/berquerant/gbrowse/blob/0123456789abcdef/sub/dir/main.go#L42",
		},
		{
			name:   "file without line number",
			target: "main.go",
			want:   "https://github.com/berquerant/gbrowse/blob/0123456789abcdef/sub/dir/main.go",
		},
		{
			name:   "directory target (exists on disk)",
			target: ".",
			want:   "https://github.com/berquerant/gbrowse/blob/0123456789abcdef/sub/dir",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			target, err := parse.ReadTarget(tc.target)
			assert.Nil(t, err)

			got, err := urlx.Build(context.Background(), gitCmd, target)
			if tc.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
