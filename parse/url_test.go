package parse_test

import (
	"testing"

	"github.com/berquerant/gbrowse/parse"
	"github.com/stretchr/testify/assert"
)

func TestReadRepoUrl(t *testing.T) {
	for _, tc := range []struct {
		title string
		value string
		want  string
	}{
		{
			title: "github",
			value: "git@github.com:berquerant/pdotfiles.git",
			want:  "https://github.com/berquerant/pdotfiles",
		},
		{
			title: "other",
			value: "git@gitlab.com:foo/group/repo.git",
			want:  "https://gitlab.com/foo/group/repo",
		},
		{
			title: "ssh+git",
			value: "ssh://git@github.com/berquerant/rpath.git",
			want:  "https://github.com/berquerant/rpath",
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got := parse.ReadRepoURL(tc.value)
			assert.Equal(t, tc.want, got)
		})
	}
}
