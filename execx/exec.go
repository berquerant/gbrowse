package execx

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/berquerant/gbrowse/ctxlog"
	"go.uber.org/zap"
)

func Run(ctx context.Context, command string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, arg...)
	cmd.Dir = "."

	var (
		stdout strings.Builder
		stderr strings.Builder
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger := ctxlog.From(ctx)
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("failed to run %s, %w: %s",
			strings.Join(cmd.Args, " "), err, stderr.String())

		logger.Debug("failed to execx.Run",
			zap.String("command", strings.Join(cmd.Args, " ")),
			zap.Strings("command_list", cmd.Args),
			zap.Error(err),
		)
		return "", err
	}

	result := strings.TrimSuffix(stdout.String(), "\n")
	logger.Debug("execx.Run",
		zap.String("command", strings.Join(cmd.Args, " ")),
		zap.Strings("command_list", cmd.Args),
		zap.String("result", result),
	)
	return result, nil
}
