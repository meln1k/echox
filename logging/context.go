package logging

import (
	"context"
	"log/slog"
)

// AppendCtx appends slog.Attrs to the context, to be used in logging handlers.
func AppendCtx(parent context.Context, attrs ...slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attrs...)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attrs...)
	return context.WithValue(parent, slogFields, v)
}
