//go:build linux

package browse

func command(target string) []string {
	return []string{
		"xdg-open", target,
	}
}
