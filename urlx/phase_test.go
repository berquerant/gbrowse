package urlx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/berquerant/gbrowse/config"
	"github.com/berquerant/gbrowse/urlx"
	"github.com/stretchr/testify/assert"
)

type mockGit struct {
	defaultBranch  string
	headObjectName string
	describeTag    string
	showCurrent    string
	commitHash     string
}

func (g *mockGit) DefaultBranch(ctx context.Context) (string, error) {
	return g.defaultBranch, nil
}
func (*mockGit) RemoteOriginURL(ctx context.Context) (string, error) {
	return "", nil
}
func (g *mockGit) HeadObjectName(ctx context.Context) (string, error) {
	return g.headObjectName, nil
}
func (*mockGit) ShowPrefix(ctx context.Context) (string, error) {
	return "", nil
}
func (*mockGit) RelativePath(ctx context.Context, path string) (string, error) {
	return "", nil
}
func (g *mockGit) DescribeTag(ctx context.Context) (string, error) {
	return g.describeTag, nil
}
func (g *mockGit) ShowCurrent(ctx context.Context) (string, error) {
	return g.showCurrent, nil
}
func (g *mockGit) CommitHash(ctx context.Context) (string, error) {
	return g.commitHash, nil
}

type mockCustomPhaseExecutor struct {
	id     string
	result string
}

func (c *mockCustomPhaseExecutor) Execute(_ context.Context, id string) (string, error) {
	if id == c.id {
		return c.result, nil
	}
	return "", errors.New("mock custom")
}

func TestPhaseExecutor(t *testing.T) {
	t.Run("builtin", func(t *testing.T) {
		const (
			customID     = "custom-id"
			customResult = "custom-result"
		)

		for _, tc := range []struct {
			title      string
			phase      config.Phase
			gitCommand *mockGit
			want       string
			err        error
		}{
			{
				title: "branch ignores HEAD",
				phase: config.NewPhase(config.Pbranch),
				gitCommand: &mockGit{
					headObjectName: "HEAD",
				},
				err: urlx.ErrGetBranch,
			},
			{
				title: "branch",
				phase: config.NewPhase(config.Pbranch),
				gitCommand: &mockGit{
					headObjectName: "mybranch",
				},
				want: "mybranch",
			},
			{
				title: "default branch",
				phase: config.NewPhase(config.PdefaultBranch),
				gitCommand: &mockGit{
					defaultBranch: "mydefault",
				},
				want: "mydefault",
			},
			{
				title: "tag but not detached HEAD",
				phase: config.NewPhase(config.Ptag),
				gitCommand: &mockGit{
					showCurrent: "master",
				},
				err: urlx.ErrGetTag,
			},
			{
				title: "tag",
				phase: config.NewPhase(config.Ptag),
				gitCommand: &mockGit{
					showCurrent: "",
					describeTag: "mytag",
				},
				want: "mytag",
			},
			{
				title: "commit",
				phase: config.NewPhase(config.Pcommit),
				gitCommand: &mockGit{
					commitHash: "mycommit",
				},
				want: "mycommit",
			},
			{
				title: "custom",
				phase: config.NewPhase(customID),
				want:  customResult,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				p := urlx.NewPhaseExecutor(tc.gitCommand, &mockCustomPhaseExecutor{
					id:     customID,
					result: customResult,
				})
				got, err := p.Execute(context.TODO(), tc.phase)
				if tc.err != nil {
					assert.ErrorIs(t, err, tc.err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})
}
