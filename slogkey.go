package errors

import (
	"log/slog"
)

// Log is used to log an error as a field with the key "err" in slog loggers.
// For some reason, despite the prevalence of this pattern in other structured loggers,
// similar functionality was not included in slog.
//
// Usage:
//
//		logger.Error("failed to get user data",
//	     slog.Str("actor-id", actor.ID),
//	     slog.Str("user-id", user.ID),
//	     errors.Log(err),
//	 )
//
// WARNING the purpose is questionable, may be to remove it?
func Log(err error) slog.Attr {
	return slog.Any("err", err)
}
