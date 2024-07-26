package skylight

import (
	"fmt"
	"slices"

	"github.com/sirupsen/logrus"
)

var defaultClient = New().WithStandardLogger()

type Level int

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "fatal"
	case LevelNone:
		return ""
	default:
		return ""
	}
}

const (
	// LevelNone represents no logging (no-op).
	LevelNone Level = iota

	// Very detailed logs, for deep troubleshooting.
	LevelTrace

	// Detailed logs useful during development.
	LevelDebug

	// High-level operational information.
	LevelInfo

	// Indications of potential issues.
	LevelWarn

	// Errors that prevent certain functions but not the whole application.
	LevelError

	// Severe errors causing application shutdown.
	LevelFatal

	// LevelPanic represents critical errors that cause the application to panic and terminate abruptly.
	LevelPanic
)

type Client struct {
	level     Level
	observers []*Observer
}

func New(opts ...Option) *Client {
	c := &Client{
		level: LevelInfo,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) WithObserver(o ...*Observer) *Client {
	if c == nil {
		return nil
	}
	c.observers = append(c.observers, o...)
	return c
}

func (c *Client) WithLevel(level Level) *Client {
	if c == nil {
		return nil
	}
	c.level = level
	return c
}

func (c *Client) WithStandardLogger() *Client {
	if c == nil {
		return nil
	}

	hasLogger := slices.ContainsFunc(c.observers, func(o *Observer) bool {
		return o.id == "logger"
	})

	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)

	if !hasLogger {
		c.WithObserver(&Observer{
			id:   "logger",
			cond: func(e *Event) bool { return true },
			handler: func(e *Event) {
				fields := map[string]any{}
				if len(e.fields) > 0 {
					fields = e.fields
				}
				if e.parentID != "" {
					fields["parentID"] = e.parentID
				}
				message := e.message
				if e.topic != "" {
					message = fmt.Sprintf("[%s] %s", e.topic, message)
				}

				switch e.level {
				case LevelTrace:
					logger.WithFields(fields).Trace(message)
				case LevelDebug:
					logger.WithFields(fields).Debug(message)
				case LevelInfo:
					logger.WithFields(fields).Info(message)
				case LevelWarn:
					logger.WithFields(fields).Warn(message)
				case LevelError:
					logger.WithFields(fields).Error(message)
				case LevelFatal:
					logger.WithFields(fields).Fatal(message)
				case LevelPanic:
					logger.WithFields(fields).Panic(message)
				}
			},
		})
	}

	return c
}

// Trace creates a new event with the trace level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func (c *Client) Trace(args ...any) *Event {
	return newEvent(LevelTrace, fmt.Sprint(args...), c)
}

// Tracef creates a new event with the trace level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func (c *Client) Tracef(v string, args ...any) *Event {
	return newEvent(LevelTrace, fmt.Sprintf(v, args...), c)
}

// Debug creates a new event with the debug level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func (c *Client) Debug(args ...any) *Event {
	return newEvent(LevelDebug, fmt.Sprint(args...), c)
}

// Debugf creates a new event with the debug level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func (c *Client) Debugf(v string, args ...any) *Event {
	return newEvent(LevelDebug, fmt.Sprintf(v, args...), c)
}

// Info creates a new event with the info level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func (c *Client) Info(args ...any) *Event {
	return newEvent(LevelInfo, fmt.Sprint(args...), c)
}

// Infof creates a new event with the info level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func (c *Client) Infof(v string, args ...any) *Event {
	return newEvent(LevelInfo, fmt.Sprintf(v, args...), c)
}

// Warn creates a new event with the warn level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func (c *Client) Warn(args ...any) *Event {
	return newEvent(LevelWarn, fmt.Sprint(args...), c)
}

// Warnf creates a new event with the warn level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func (c *Client) Warnf(v string, args ...any) *Event {
	return newEvent(LevelWarn, fmt.Sprintf(v, args...), c)
}

// Error creates a new event with the error level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func (c *Client) Error(args ...any) *Event {
	return newEvent(LevelError, fmt.Sprint(args...), c)
}

// Errorf creates a new event with the error level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func (c *Client) Errorf(v string, args ...any) *Event {
	return newEvent(LevelError, fmt.Sprintf(v, args...), c)
}

// Fatal creates a new event with the fatal level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func (c *Client) Fatal(args ...any) *Event {
	return newEvent(LevelFatal, fmt.Sprint(args...), c)
}

// Fatalf creates a new event with the fatal level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func (c *Client) Fatalf(v string, args ...any) *Event {
	return newEvent(LevelFatal, fmt.Sprintf(v, args...), c)
}

// Panic creates a new event with the panic level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func (c *Client) Panic(args ...any) *Event {
	return newEvent(LevelPanic, fmt.Sprint(args...), c)
}

// Panicf creates a new event with the panic level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func (c *Client) Panicf(v string, args ...any) *Event {
	return newEvent(LevelPanic, fmt.Sprintf(v, args...), c)
}

// Trace creates a new event with the trace level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func Trace(args ...any) *Event {
	return newEvent(LevelTrace, fmt.Sprint(args...), defaultClient)
}

// Tracef creates a new event with the trace level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func Tracef(v string, args ...any) *Event {
	return newEvent(LevelTrace, fmt.Sprintf(v, args...), defaultClient)
}

// Debug creates a new event with the debug level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func Debug(args ...any) *Event {
	return newEvent(LevelDebug, fmt.Sprint(args...), defaultClient)
}

// Debugf creates a new event with the debug level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func Debugf(v string, args ...any) *Event {
	return newEvent(LevelDebug, fmt.Sprintf(v, args...), defaultClient)
}

// Info creates a new event with the info level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func Info(args ...any) *Event {
	return newEvent(LevelInfo, fmt.Sprint(args...), defaultClient)
}

// Infof creates a new event with the info level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func Infof(v string, args ...any) *Event {
	return newEvent(LevelInfo, fmt.Sprintf(v, args...), defaultClient)
}

// Warn creates a new event with the warn level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func Warn(args ...any) *Event {
	return newEvent(LevelWarn, fmt.Sprint(args...), defaultClient)
}

// Warnf creates a new event with the warn level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func Warnf(v string, args ...any) *Event {
	return newEvent(LevelWarn, fmt.Sprintf(v, args...), defaultClient)
}

// Error creates a new event with the error level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func Error(args ...any) *Event {
	return newEvent(LevelError, fmt.Sprint(args...), defaultClient)
}

// Errorf creates a new event with the error level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func Errorf(v string, args ...any) *Event {
	return newEvent(LevelError, fmt.Sprintf(v, args...), defaultClient)
}

// Fatal creates a new event with the fatal level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func Fatal(args ...any) *Event {
	return newEvent(LevelFatal, fmt.Sprint(args...), defaultClient)
}

// Fatalf creates a new event with the fatal level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func Fatalf(v string, args ...any) *Event {
	return newEvent(LevelFatal, fmt.Sprintf(v, args...), defaultClient)
}

// Panic creates a new event with the panic level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new event.
func Panic(args ...any) *Event {
	return newEvent(LevelPanic, fmt.Sprint(args...), defaultClient)
}

// Panicf creates a new event with the panic level and associates it with the client.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new event.
func Panicf(v string, args ...any) *Event {
	return newEvent(LevelPanic, fmt.Sprintf(v, args...), defaultClient)
}

func Default() *Client {
	return defaultClient
}
