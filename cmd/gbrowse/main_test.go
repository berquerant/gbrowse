package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EnvMap struct {
	DefaultBranch   string `json:"default_branch"`
	RemoteOriginURL string `json:"remote_origin_url"`
	HeadObjectName  string `json:"head_object_name"`
	ShowPrefix      string `json:"show_prefix"`
	RelativePath    string `json:"relative_path"`
	DescribeTag     string `json:"describe_tag"`
	ShowCurrent     string `json:"show_current"`
	CommitHash      string `json:"commit_hash"`
}

func (e *EnvMap) JSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func defaultEnvMap() *EnvMap {
	return &EnvMap{
		DefaultBranch:   "master",
		RemoteOriginURL: "remote-origin",
		HeadObjectName:  "head-object",
		ShowPrefix:      "show-prefix",
		RelativePath:    "relative-path",
		DescribeTag:     "describe-tag",
		ShowCurrent:     "show-current",
		CommitHash:      "commit-hash",
	}
}

func TestEndToEnd(t *testing.T) {
	e := newExecutor(t)
	defer e.close()
	var (
		newEnvSlices = func(envMap *EnvMap) []string {
			return []string{
				fmt.Sprintf("GIT=%s", e.git),                        // mock git binary
				fmt.Sprintf("GBROWSE_GIT_CONFIG=%s", envMap.JSON()), // see cmd/gbrowse-git
			}
		}
	)

	t.Run("help", func(t *testing.T) {
		envSlices := newEnvSlices(defaultEnvMap())
		output, err := run(envSlices, e.cmd, "-h")
		assert.Nil(t, err)
		t.Log(string(output))
	})

	t.Run("run", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			envs := defaultEnvMap()
			envSlices := newEnvSlices(envs)

			for _, tc := range []struct {
				name string
				opt  []string
				want string
			}{
				{
					name: "root",
					opt:  []string{"-print"},
					want: strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix}, "/"),
				},
				{
					name: "dir",
					opt:  []string{"-print", "dir"},
					want: strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					name: "dir/file",
					opt:  []string{"-print", "dir/file"},
					want: strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix, "dir/file"}, "/"),
				},
				{
					name: "linum",
					opt:  []string{"-print", "dir/file:10"},
					want: strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix, "dir/file#L10"}, "/"),
				},
			} {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					output, err := run(envSlices, e.cmd, tc.opt...)
					assert.Nil(t, err)
					assert.Equal(t, tc.want, string(output))
				})
			}
		})
	})

}

func run(env []string, name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	cmd.Env = env
	cmd.Dir = "."
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

type executor struct {
	dir string
	cmd string
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
	cmd := filepath.Join(dir, "gbrowse")
	// build gbrowse command
	if _, err := run(nil, "go", "build", "-o", cmd); err != nil {
		t.Fatal(err)
	}
	e.dir = dir
	e.cmd = cmd

	git := filepath.Join(dir, "gbrowse-git")
	// build gbrowse-git command
	if _, err := run(nil, "go", "build", "-o", git, "../gbrowse-git"); err != nil {
		t.Fatal(err)
	}
	e.git = git
}

func (e *executor) close() {
	os.RemoveAll(e.dir)
}
