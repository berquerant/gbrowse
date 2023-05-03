package main

import (
	"fmt"
	"strings"
)

var (
	AuthorName  string
	AuthorEmail string
	Version     string
	GoVersion   string
	Commit      string
	Project     string
	GithubUser  string
)

func Ldflags() string {
	flags := []struct {
		key   string
		value string
	}{
		{
			key:   "Author",
			value: AuthorName,
		},
		{
			key:   "Email",
			value: AuthorEmail,
		},
		{
			key:   "Repository",
			value: fmt.Sprintf("https://github.com/%s/%s", GithubUser, Project),
		},
		{
			key:   "Version",
			value: Version,
		},
		{
			key:   "Go Version",
			value: GoVersion,
		},
		{
			key:   "Commit",
			value: Commit,
		},
	}

	var ss []string
	for _, f := range flags {
		if f.value != "" {
			ss = append(ss, fmt.Sprintf("%-12s:\t%s", f.key, f.value))
		}
	}
	return strings.Join(ss, "\n")
}
