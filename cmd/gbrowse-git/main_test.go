package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	e := newExecutor(t)
	defer e.close()

	const (
		defaultBranch   = "master"
		remoteOriginURL = "remote-origin"
		headObjectName  = "head-object"
		showPrefix      = "show-prefix"
		relativePath    = "relative-path"
		describeTag     = "describe-tag"
		showCurrent     = "show-current"
		commitHash      = "commit-hash"
	)
	envBytes, _ := json.Marshal(map[string]string{
		"default_branch":    defaultBranch,
		"remote_origin_url": remoteOriginURL,
		"head_object_name":  headObjectName,
		"show_prefix":       showPrefix,
		"relative_path":     relativePath,
		"describe_tag":      describeTag,
		"show_current":      showCurrent,
		"commit_hash":       commitHash,
	})
	envSlices := []string{
		fmt.Sprintf("GBROWSE_GIT=%s", e.git),
		fmt.Sprintf("GBROWSE_GIT_CONFIG=%s", string(envBytes)),
	}

	t.Run("gbrowse-git", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			args []string
			want string
		}{
			{
				name: "DefaultBranch",
				args: []string{"remote", "show", "origin"},
				want: fmt.Sprintf("HEAD branch: %s", defaultBranch),
			},
			{
				name: "RemoteOriginURL",
				args: []string{"config", "--get", "remote.origin.url"},
				want: remoteOriginURL,
			},
			{
				name: "HeadObjectName",
				args: []string{"rev-parse", "--abbrev-ref", "@"},
				want: headObjectName,
			},
			{
				name: "ShowPrefix",
				args: []string{"rev-parse", "--show-prefix"},
				want: showPrefix,
			},
			{
				name: "RelativePath",
				args: []string{"ls-files", "--full-name"},
				want: relativePath,
			},
			{
				name: "DescribeTag",
				args: []string{"describe", "--tags", "--abbrev=0"},
				want: describeTag,
			},
			{
				name: "ShowCurrent",
				args: []string{"branch", "--show-current"},
				want: showCurrent,
			},
			{
				name: "CommitHash",
				args: []string{"rev-parse", "@"},
				want: commitHash,
			},
		} {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				output, err := run(envSlices, e.git, tc.args...)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, string(output))
			})
		}
	})
}

func run(env []string, name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	cmd.Env = env
	cmd.Dir = "."
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

type executor struct {
	dir string
	git string
}

func newExecutor(t *testing.T) *executor {
	t.Helper()
	e := &executor{}
	e.init(t)
	return e
}

func (e *executor) init(t *testing.T) {
	t.Helper()
	dir, err := os.MkdirTemp("", "gbrowse")
	if err != nil {
		t.Fatal(err)
	}

	git := filepath.Join(dir, "gbrowse-git")
	// build gbrowse-git command
	if _, err := run(nil, "go", "build", "-o", git); err != nil {
		t.Fatal(err)
	}
	e.dir = dir
	e.git = git
}

func (e *executor) close() {
	os.RemoveAll(e.dir)
}
