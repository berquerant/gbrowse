package urlx

import (
	"context"
	"errors"
	"fmt"

	"github.com/berquerant/gbrowse/config"
	"github.com/berquerant/gbrowse/ctxlog"
	"github.com/berquerant/gbrowse/git"
)

func ExecutePhases(ctx context.Context, phases []config.Phase, executor PhaseExecutor) (string, error) {
	var (
		// fallback to commit hash
		actualPhases = append(phases, config.NewPhase(config.Pcommit))
		retErr       error
		logger       = ctxlog.From(ctx)
	)
	for i, p := range actualPhases {
		logger.Debug("phase", ctxlog.I("index", i), ctxlog.S("phase", p.String()))
		r, err := executor.Execute(ctx, p)
		if err == nil {
			return r, nil
		}
		logger.Debug("phase", ctxlog.I("index", i), ctxlog.Err(err))
		retErr = errors.Join(retErr, fmt.Errorf("%w: phase[%d] %s", err, i, p))
	}
	return "", retErr
}

type PhaseExecutor interface {
	Execute(ctx context.Context, phase config.Phase) (string, error)
}

func NewPhaseExecutor(gitCommand git.Git, customExecutor CustomPhaseExecutor) PhaseExecutor {
	return &phaseExecutor{
		gitCommand:     gitCommand,
		customExecutor: customExecutor,
	}
}

type phaseExecutor struct {
	gitCommand     git.Git
	customExecutor CustomPhaseExecutor
}

var (
	ErrGetBranch      = errors.New("GetBranch")
	ErrGetTag         = errors.New("GetTag")
	ErrUnknownBuiltin = errors.New("UnknownBuiltin")
)

func (e *phaseExecutor) Execute(ctx context.Context, phase config.Phase) (string, error) {
	r, err := e.executeBuiltin(ctx, phase)
	if errors.Is(err, ErrUnknownBuiltin) {
		return e.executeCustom(ctx, phase)
	}
	return r, err
}

func (e *phaseExecutor) executeBuiltin(ctx context.Context, phase config.Phase) (string, error) {
	switch phase.String() {
	case config.Pbranch:
		r, err := e.gitCommand.HeadObjectName(ctx)
		if err != nil {
			return "", err
		}
		if r == "HEAD" {
			return "", fmt.Errorf("%w: got HEAD", ErrGetBranch)
		}
		return r, nil
	case config.PdefaultBranch:
		return e.gitCommand.DefaultBranch(ctx)
	case config.Ptag:
		if r, _ := e.gitCommand.ShowCurrent(ctx); r != "" {
			return "", fmt.Errorf("%w: not `detached HEAD` state", ErrGetTag)
		}
		return e.gitCommand.DescribeTag(ctx)
	case config.Pcommit:
		return e.gitCommand.CommitHash(ctx)
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownBuiltin, phase)
	}
}

func (e *phaseExecutor) executeCustom(ctx context.Context, phase config.Phase) (string, error) {
	return e.customExecutor.Execute(ctx, phase.String())
}
