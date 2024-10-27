package echox

import (
	"context"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/meln1k/echox/sse"
)

// Render renders templates with 200, most commonly used for htmx
func Render(c echo.Context, templates ...templ.Component) error {
	return RenderStatus(c, 200, templates...)
}

// RenderStatus renders templates to the response with the given status code
func RenderStatus(c echo.Context, statusCode int, templates ...templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := RenderW(c, buf, templates...); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}

type contextKey string

const (
	reverseContextKey contextKey = "reverse"
)

// RenderW renders templates to the given writer, and adds the reverse function to the template context.
func RenderW(c echo.Context, writer io.Writer, templates ...templ.Component) error {
	ctx := c.Request().Context()

	ctx = context.WithValue(ctx, reverseContextKey, c.Echo().Reverse)
	for _, t := range templates {
		if err := t.Render(ctx, writer); err != nil {
			return err
		}
	}
	return nil
}

// ReverseUrl returns the URL for the given handler name and params.
func ReverseUrl(ctx context.Context, handlerName string, params ...any) templ.SafeURL {
	reverse, ok := ctx.Value(reverseContextKey).(func(string, ...any) string)
	if !ok {
		panic("can't get reverse function from the context, did you forget to call Render on the template?")
	}

	return templ.URL(reverse(handlerName, params...))
}

// ReverseHx returns the URL for the given handler name and params, used for HTMX templ attributes.
func ReverseHx(c context.Context, handlerName string, params ...any) string {
	return string(ReverseUrl(c, handlerName, params...))
}

type ServerSentEvent struct {
	Event string
	Data  templ.Component
}

func SSE(c echo.Context, eventStream <-chan ServerSentEvent) error {
	ctx := c.Request().Context()

	w := c.Response()

	w.Header().Set(echo.HeaderContentType, "text/event-stream")
	w.Header().Set(echo.HeaderCacheControl, "no-cache")
	w.Header().Set(echo.HeaderConnection, "keep-alive")

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-eventStream:
			// check if the channel has been closed
			if !ok {
				return nil
			}

			sb := strings.Builder{}
			if err := RenderW(c, &sb, event.Data); err != nil {
				return err
			}

			sseEvt := sse.Event{
				Event: []byte(event.Event),
				Data:  []byte(sb.String()),
			}

			// write the event to the client
			if err := sseEvt.MarshalTo(w); err != nil {
				return err
			}

			w.Flush()
		}
	}

}
