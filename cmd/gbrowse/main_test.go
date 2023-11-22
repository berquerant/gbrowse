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
				fmt.Sprintf("GBROWSE_GIT=%s", e.git),                // mock git binary
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
		tmpDir := t.TempDir()

		t.Run("config custom", func(t *testing.T) {
			envs := defaultEnvMap()
			envSlices := newEnvSlices(envs)
			// run gbrowse-git branch --show-current as custom phase
			configTemplate := `{"phases":["custom"],"defs":[{"id":"custom","cmd":["%s","branch","--show-current"]}]}`
			config := fmt.Sprintf(configTemplate, e.git)
			want := strings.Join([]string{envs.RemoteOriginURL, "blob", envs.ShowCurrent, envs.ShowPrefix}, "/")
			output, err := run(envSlices, e.cmd, "-print", "-config", config)
			assert.Nil(t, err)
			assert.Equal(t, want, string(output))
		})

		t.Run("config flag overwride env", func(t *testing.T) {
			envs := defaultEnvMap()
			envSlices := append(
				newEnvSlices(envs),
				`GBROWSE_CONFIG={"phases":["commit"]}`,
			)

			want := strings.Join([]string{envs.RemoteOriginURL, "blob", envs.DefaultBranch, envs.ShowPrefix}, "/")
			output, err := run(envSlices, e.cmd, "-print", "-config", `{"phases":["default_branch"]}`)
			assert.Nil(t, err)
			assert.Equal(t, want, string(output))
		})

		t.Run("config env", func(t *testing.T) {
			envs := defaultEnvMap()
			envSlices := append(
				newEnvSlices(envs),
				`GBROWSE_CONFIG={"phases":["default_branch"]}`,
			)
			want := strings.Join([]string{envs.RemoteOriginURL, "blob", envs.DefaultBranch, envs.ShowPrefix}, "/")
			output, err := run(envSlices, e.cmd, "-print")
			assert.Nil(t, err)
			assert.Equal(t, want, string(output))
		})

		t.Run("config", func(t *testing.T) {
			envs := defaultEnvMap()

			for _, tc := range []struct {
				title  string
				envs   func(*EnvMap)
				config string
				opt    []string
				want   string
			}{
				{
					title:  "empty",
					config: `{}`,
					opt:    []string{"dir"},
					want:   strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title:  "commit",
					config: `{"phases":["commit"]}`,
					opt:    []string{"dir"},
					want:   strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title:  "default branch",
					config: `{"phases":["default_branch"]}`,
					opt:    []string{"dir"},
					want:   strings.Join([]string{envs.RemoteOriginURL, "blob", envs.DefaultBranch, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title: "default branch by -phase",
					opt:   []string{"-phase", "default_branch", "dir"},
					want:  strings.Join([]string{envs.RemoteOriginURL, "blob", envs.DefaultBranch, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title:  "branch ignore HEAD",
					config: `{"phases":["branch"]}`,
					envs: func(e *EnvMap) {
						e.HeadObjectName = "HEAD"
					},
					opt:  []string{"dir"},
					want: strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title:  "branch",
					config: `{"phases":["branch"]}`,
					opt:    []string{"dir"},
					want:   strings.Join([]string{envs.RemoteOriginURL, "blob", envs.HeadObjectName, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title:  "tag ignored wen not detacjed HEAD",
					config: `{"phases":["tag"]}`,
					opt:    []string{"dir"},
					want:   strings.Join([]string{envs.RemoteOriginURL, "blob", envs.CommitHash, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title:  "tag",
					config: `{"phases":["tag"]}`,
					envs: func(e *EnvMap) {
						e.ShowCurrent = ""
					},
					opt:  []string{"dir"},
					want: strings.Join([]string{envs.RemoteOriginURL, "blob", envs.DescribeTag, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title:  "fallback",
					config: `{"phases":["tag","default_branch"]}`,
					opt:    []string{"dir"},
					want:   strings.Join([]string{envs.RemoteOriginURL, "blob", envs.DefaultBranch, envs.ShowPrefix, "dir"}, "/"),
				},
				{
					title: "fallback by -phase",
					opt:   []string{"-phase", "tag,default_branch", "dir"},
					want:  strings.Join([]string{envs.RemoteOriginURL, "blob", envs.DefaultBranch, envs.ShowPrefix, "dir"}, "/"),
				},
			} {
				tc := tc
				t.Run(tc.title, func(t *testing.T) {
					envs := defaultEnvMap()
					if tc.envs != nil {
						tc.envs(envs)
					}
					envSlices := newEnvSlices(envs)
					configFile, err := os.CreateTemp(tmpDir, "")
					if err != nil {
						t.Fatal(err)
					}
					defer configFile.Close()
					if _, err := fmt.Fprint(configFile, tc.config); err != nil {
						t.Fatal(err)
					}

					opt := append([]string{"-print", "-config", configFile.Name()}, tc.opt...)
					output, err := run(envSlices, e.cmd, opt...)
					assert.Nil(t, err)
					assert.Equal(t, tc.want, string(output))
				})
			}
		})

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
