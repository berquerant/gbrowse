package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/berquerant/gbrowse/config"
	"github.com/stretchr/testify/assert"
)

func TestParseStringOrFile(t *testing.T) {
	t.Run("empty return default", func(t *testing.T) {
		assert.Equal(t, config.Default(), config.ParseStringOrFile("", nil))
	})

	t.Run("neither file nor config", func(t *testing.T) {
		assert.Equal(t, config.Default(), config.ParseStringOrFile("{", nil))
	})

	t.Run("config", func(t *testing.T) {
		want := config.Default()
		want.Phases = []config.Phase{
			config.NewPhase(config.Pcommit),
		}
		assert.Equal(t, want, config.ParseStringOrFile(`{"phases":["commit"]}`, nil))
	})

	t.Run("file", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if _, err := fmt.Fprint(f, `{"phases":["commit"]}`); err != nil {
			t.Fatal(err)
		}

		want := config.Default()
		want.Phases = []config.Phase{
			config.NewPhase(config.Pcommit),
		}
		assert.Equal(t, want, config.ParseStringOrFile(f.Name(), nil))
	})

	t.Run("overwrite by config", func(t *testing.T) {
		origin := config.Default()
		origin.Phases = []config.Phase{
			config.NewPhase(config.Ptag),
		}
		want := config.Default()
		want.Phases = []config.Phase{
			config.NewPhase(config.Pcommit),
		}
		assert.Equal(t, want, config.ParseStringOrFile(`{"phases":["commit"]}`, origin))
	})

	t.Run("phases and defs", func(t *testing.T) {
		want := config.Default()
		want.Phases = []config.Phase{
			config.NewPhase("custom"),
		}
		want.Definitions = []config.Definition{
			{
				ID:      "custom",
				Command: []string{"pwd"},
			},
		}
		assert.Equal(t, want, config.ParseStringOrFile(`{"phases":["custom"],"defs":[{"id":"custom","cmd":["pwd"]}]}`, nil))
	})
}
