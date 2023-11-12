package urlx

import (
	"context"
	"errors"
	"fmt"

	"github.com/berquerant/gbrowse/config"
	"github.com/berquerant/gbrowse/execx"
)

var (
	ErrDefinitionNotFound = errors.New("DefinitionNotFound")
)

type CustomPhaseExecutor interface {
	Execute(ctx context.Context, id string) (string, error)
}

func NewCustomPhaseExecutor(definitions []config.Definition) CustomPhaseExecutor {
	d := map[string]config.Definition{}
	for _, def := range definitions {
		d[def.ID] = def
	}
	return &customPhaseExecutor{
		definitions: d,
	}
}

type customPhaseExecutor struct {
	definitions map[string]config.Definition // id to defs
}

func (e *customPhaseExecutor) Execute(ctx context.Context, id string) (string, error) {
	def, ok := e.definitions[id]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrDefinitionNotFound, id)
	}

	return execx.Run(ctx, def.Command[0], def.Command[1:]...)
}
