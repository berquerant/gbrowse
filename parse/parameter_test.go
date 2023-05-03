package parse_test

import (
	"testing"

	"github.com/berquerant/gbrowse/parse"
	"github.com/stretchr/testify/assert"
)

func TestReadTarget(t *testing.T) {
	t.Run("invalid target", func(t *testing.T) {
		_, err := parse.ReadTarget("a:b")
		assert.NotNil(t, err)
	})

	for _, tc := range []struct {
		title string
		value string
		want  *parse.Target
	}{
		{
			title: "empty",
			value: "",
			want:  parse.NewPathTarget(""),
		},
		{
			title: "path only",
			value: "a",
			want:  parse.NewPathTarget("a"),
		},
		{
			title: "path and linum",
			value: "a:1",
			want:  parse.NewTarget("a", 1),
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, err := parse.ReadTarget(tc.value)
			assert.Nil(t, err)
			assert.Equal(t, tc.want.Path(), got.Path())

			wantLinum, wantHasLinum := tc.want.Linum()
			gotLinum, gotHasLinum := got.Linum()
			assert.Equal(t, wantHasLinum, gotHasLinum)
			if wantHasLinum {
				assert.Equal(t, wantLinum, gotLinum)
			}
		})
	}
}
