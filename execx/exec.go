package execx

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/berquerant/gbrowse/ctxlog"
)

func Run(ctx context.Context, command string, arg ...string) (string, error) {
	var (
		logger = ctxlog.From(ctx)
		cmd    = exec.CommandContext(ctx, command, arg...)
		genErr = func(err error) error {
			err = fmt.Errorf("failed to run %s, %w", strings.Join(cmd.Args, " "), err)
			logger.Debug("failed to execx.Run",
				ctxlog.S("command", strings.Join(cmd.Args, " ")),
				ctxlog.SS("command_list", cmd.Args),
				ctxlog.Err(err),
			)
			return err
		}
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd.Env = os.Environ()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", genErr(fmt.Errorf("%w: %s", err, stderr.String()))
	}

	result := strings.TrimSuffix(stdout.String(), "\n")
	logger.Debug("execx.Run",
		ctxlog.S("command", strings.Join(cmd.Args, " ")),
		ctxlog.SS("command_list", cmd.Args),
		ctxlog.S("result", result),
	)
	return result, nil
}
