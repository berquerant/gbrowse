package urlx_test

import (
	"context"
	"testing"

	"github.com/berquerant/gbrowse/config"
	"github.com/berquerant/gbrowse/urlx"
	"github.com/stretchr/testify/assert"
)

func TestCustomPhaseExecutor(t *testing.T) {
	var (
		defID       = "echo"
		defResult   = "custom-phase"
		definitions = []config.Definition{
			{
				ID:      defID,
				Command: []string{"echo", defResult},
			},
		}
		executor = urlx.NewCustomPhaseExecutor(definitions)
	)

	t.Run("not found", func(t *testing.T) {
		_, err := executor.Execute(context.TODO(), "not found")
		assert.ErrorIs(t, err, urlx.ErrDefinitionNotFound)
	})

	t.Run("found", func(t *testing.T) {
		got, err := executor.Execute(context.TODO(), defID)
		assert.Nil(t, err)
		assert.Equal(t, defResult, got)
	})
}
