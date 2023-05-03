package browse

import (
	"context"
	"fmt"

	"github.com/berquerant/gbrowse/execx"
)

// Run opens the target in the browser.
func Run(ctx context.Context, target string) error {
	cmd := command(target)
	if _, err := execx.Run(ctx, cmd[0], cmd[1:]...); err != nil {
		return fmt.Errorf("failed to browse: %w", err)
	}
	return nil
}
