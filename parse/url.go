package parse

import "strings"

func ReadRepoUrl(value string) string {
	replaceTuples := []struct {
		from, to string
	}{
		{
			from: ":",
			to:   "/",
		},
		{
			from: "git@",
			to:   "https///",
		},
		{
			from: "git///",
			to:   "https///",
		},
		{
			from: ".git",
			to:   "",
		},
		{
			from: "https///",
			to:   "https://",
		},
	}

	for _, t := range replaceTuples {
		value = strings.ReplaceAll(value, t.from, t.to)
	}
	return value
}
