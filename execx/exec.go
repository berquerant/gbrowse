package execx

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/berquerant/execx"
	"github.com/berquerant/gbrowse/ctxlog"
)

func Run(ctx context.Context, command string, arg ...string) (result string, retErr error) {
	var (
		logger = ctxlog.From(ctx)
		cmd    = execx.New(command, arg...)
		genErr = func(err error) error {
			err = fmt.Errorf("failed to run %s, %w", strings.Join(cmd.Args, " "), err)
			logger.Debug("failed to execx.Run",
				ctxlog.S("command", strings.Join(cmd.Args, " ")),
				ctxlog.SS("command_list", cmd.Args),
				ctxlog.Err(err),
			)
			return err
		}
	)

	cmd.Env = execx.EnvFromEnviron()
	r, err := cmd.Run(ctx)
	if err != nil {
		retErr = genErr(err)
		return
	}

	stdout, err := io.ReadAll(r.Stdout)
	if err != nil {
		retErr = genErr(err)
		return
	}

	result = strings.TrimSuffix(string(stdout), "\n")
	logger.Debug("execx.Run",
		ctxlog.S("command", strings.Join(cmd.Args, " ")),
		ctxlog.SS("command_list", cmd.Args),
		ctxlog.S("result", result),
	)
	return
}
