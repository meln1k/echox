package logging

import (
	"context"
	"log/slog"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

// ContextHandler is a slog.Handler that appends slog.Attrs to the context, to be used in logging handlers.
type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}
