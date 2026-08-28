package ctxlog_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/berquerant/gbrowse/ctxlog"
	"github.com/stretchr/testify/assert"
)

func TestFrom(t *testing.T) {
	orig := ctxlog.New(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	for _, tc := range []struct {
		name       string
		ctx        context.Context
		wantLogger ctxlog.Logger
	}{
		{
			name:       "nil context returns noop",
			ctx:        nil,
			wantLogger: nil,
		},
		{
			name:       "empty context returns noop",
			ctx:        context.Background(),
			wantLogger: nil,
		},
		{
			name:       "context with logger returns original logger",
			ctx:        ctxlog.With(context.Background(), orig),
			wantLogger: orig,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := ctxlog.From(tc.ctx)
			assert.NotNil(t, l)
			if tc.wantLogger != nil {
				assert.Equal(t, tc.wantLogger, l)
			}
			// ensure calling methods on returned logger never panics
			l.Debug("debug test", ctxlog.S("k", "v"))
			l.Info("info test", ctxlog.I("k", 1))
			l.Error("error test", ctxlog.B("k", true))
		})
	}
}
