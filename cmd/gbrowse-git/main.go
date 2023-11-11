package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/berquerant/gbrowse/env"
)

type config struct {
	DefaultBranch   string `json:"default_branch"`
	RemoteOriginURL string `json:"remote_origin_url"`
	HeadObjectName  string `json:"head_object_name"`
	ShowPrefix      string `json:"show_prefix"`
	RelativePath    string `json:"relative_path"`
	DescribeTag     string `json:"describe_tag"`
	ShowCurrent     string `json:"show_current"`
	CommitHash      string `json:"commit_hash"`
}

func (c *config) intoMappingTuples() mappingTupleList {
	return []*mappingTuple{
		newMappingTuple([]string{"remote", "show", "origin"}, fmt.Sprintf("HEAD branch: %s", c.DefaultBranch)),
		newMappingTuple([]string{"config", "--get", "remote.origin.url"}, c.RemoteOriginURL),
		newMappingTuple([]string{"rev-parse", "--abbrev-ref", "@"}, c.HeadObjectName),
		newMappingTuple([]string{"rev-parse", "--show-prefix"}, c.ShowPrefix),
		newMappingTuple([]string{"ls-files", "--full-name"}, c.RelativePath),
		newMappingTuple([]string{"describe", "--tags", "--abbrev=0"}, c.DescribeTag),
		newMappingTuple([]string{"branch", "--show-current"}, c.ShowCurrent),
		newMappingTuple([]string{"rev-parse", "@"}, c.CommitHash),
	}
}

type mappingTuple struct {
	key   string
	value string
}

func newMappingTuple(key []string, value string) *mappingTuple {
	return &mappingTuple{
		key:   strings.Join(key, "|"),
		value: value,
	}
}

type mappingTupleList []*mappingTuple

func (l mappingTupleList) get(args []string) (string, bool) {
	key := strings.Join(args, "|")
	for _, t := range l {
		if strings.HasPrefix(key, t.key) {
			return t.value, true
		}
	}
	return "", false
}

func (l mappingTupleList) String() string {
	m := make(map[string]string, len(l))
	for _, t := range l {
		m[t.key] = t.value
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		cValue = env.GetOr("GBROWSE_GIT_CONFIG", "{}")
		args   = os.Args[1:]
		c      config
	)

	fail(json.Unmarshal([]byte(cValue), &c))
	m := c.intoMappingTuples()
	if value, ok := m.get(args); ok {
		fmt.Print(value)
		return
	}

	fmt.Fprintf(os.Stderr, "not found, config=%s, args=%v, mapping=%s\n", cValue, args, m)
	fmt.Fprintln(os.Stderr, usage)
	os.Exit(1)
}

const usage = `gbrowse-git -- Mock git command for command.Git

Usage:
  gbrowse-git [flags]

Environment variables:
  GBROWSE_GIT_CONFIG
    mock values setting as json`
