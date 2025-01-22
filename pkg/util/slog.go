package util

import "log/slog"

func SlErr(err error) slog.Attr {
	return slog.String("error", err.Error())
}
