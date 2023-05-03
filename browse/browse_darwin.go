//go:build darwin

package browse

func command(target string) []string {
	return []string{
		"open", target,
	}
}
