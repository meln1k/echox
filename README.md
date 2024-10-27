# echox

High power tools for echo + templ + HTMX

This is just a bunch of echo/templ/HTMX helpers I found useful when making pet projects.


## Installation

```sh
go get github.com/meln1k/echox
```

## Usage

### Rendering templ components in echo handlers

Most of the time you will render templ components in echo handlers with default 200 status code.

```go
import "github.com/meln1k/echox"


templ HelloWorld() {
    <div>Hello, World!</div>
}

func handler(c echo.Context) error {
	return echox.Render(c, HelloWorld())
}
```


### Reversing URLs in templates

To use HTMX attributes like `hx-get` or `hx-post` we need to reverse URLs in templates.

```go

import "github.com/meln1k/echox/hx"

templ HelloWorld() {
    <div hx-get={ echox.ReverseHx(ctx, "handler") }>Hello, World!</div>

    <a href={ echox.ReverseUrl(ctx, "handler") }>Hello, World!</a>
}

func handler(c echo.Context) error {
	return echox.Render(c, HelloWorld())
}

echo.GET("/", handler).Name = "handler"
```


### HTMX headers, triggers, redirects, etc.

Sometimes there is a dire need to trigger some client-side HTMX things, like redirects or events.

```go
import "github.com/meln1k/echox/hx"

// Redirecting via HTMX or standard redirect, depending on the request
func handler(c echo.Context) error {
    hx.Redirect(c, "/")
}

// Triggering an event on the client
func handler(c echo.Context) error {
    hx.Trigger(c, "my-event")
}

// Triggering an event after swap (could be used for toast notifications, etc.)
func handler(c echo.Context) error {
    hx.TriggerAfterSwap(c, "my-event")
}

// for more good stuff, see the hx/htmxtools.go file

```

### Server Sent Events with htmx and templ

An example how you can send async updates to the client using HTMX and templ.


```go
import "github.com/meln1k/echox"


templ MessagesPage() {
    <div>
        <h1>Messages</h1>
        <div hx-ext="sse" sse-connect={ view.ReverseHx(ctx, "sseHandler") } hx-swap="beforeEnd">
        </div>
    </div>
}

templ Message(content string) {
    <div>{ content }</div>
}


func handler(c echo.Context) error {
    return echox.Render(c, MessagesPage())
}

func sseHandler(c echo.Context) error {

    ctx := c.Request().Context()

    events := make(chan echox.ServerSentEvent, 1)
    defer close(events)

    go func() {

        messageNumber := 0
        for  {
            time.Sleep(time.Second)

            // check if the request has been cancelled
            select {
            case <-ctx.Done():
                return
            default:
            }

            // push the new message to the event stream
            template := Message(fmt.Sprintf("Hello, World! #%d", messageNumber))
            events <- echox.ServerSentEvent{
                Event: "message",
                Data:  template,
            }

            messageNumber++
        }
    }()

    return echox.SSE(c, events)
}


```

### Slog logging with echo

While this is not htmx or templ related, it's still useful to have a good logging setup for your echo app.

```go
import "github.com/meln1k/echox/logging"

import "log/slog"

import "github.com/labstack/echo/v4"
import "github.com/labstack/echo/v4/middleware"


func main() {

    slogHandler := slog.NewSlogHandler(slog.Default())

    e := echo.New()

    e.Use(middleware.RequestLoggerWithConfig(logging.RequestLoggerConfig(slogHandler)))

    ...

```
