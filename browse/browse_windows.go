//go:build windows

package browse

func command(target string) []string {
	return []string{
		"start", target,
	}
}
