package parse_test

import (
	"testing"

	"github.com/berquerant/gbrowse/parse"
	"github.com/stretchr/testify/assert"
)

func TestReadTarget(t *testing.T) {
	for _, tc := range []struct {
		title   string
		value   string
		want    *parse.Target
		wantErr bool
	}{
		{
			title:   "invalid target",
			value:   "a:b",
			wantErr: true,
		},
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
		t.Run(tc.title, func(t *testing.T) {
			got, err := parse.ReadTarget(tc.value)
			if tc.wantErr {
				assert.NotNil(t, err)
				return
			}
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
