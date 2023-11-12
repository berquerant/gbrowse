package ctxlog

import "golang.org/x/exp/slog"

type Attr slog.Attr

func S(key, value string) Attr {
	return Attr(slog.String(key, value))
}

func B(key string, value bool) Attr {
	return Attr(slog.Bool(key, value))
}

func Any(key string, value any) Attr {
	return Attr(slog.Any(key, value))
}

func SS(key string, value []string) Attr {
	return Attr(slog.Any(key, value))
}

func Err(err error) Attr {
	return Attr(slog.Any("error", err))
}

func I(key string, value int) Attr {
	return Attr(slog.Int(key, value))
}
