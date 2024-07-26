package skylight

import (
	"fmt"
	"sync"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	eventPool = &sync.Pool{
		New: func() interface{} {
			return &Event{}
		},
	}
)

type Fields map[string]any

type Event struct {
	c         *Client
	closed    bool
	id        string
	createdAt time.Time
	emittedAt time.Time
	level     Level
	message   string
	topic     string
	parentID  string
	fields    Fields
}

func newEvent(level Level, msg string, c *Client) *Event {
	if c == nil || level < c.level {
		return nil
	}

	e := eventPool.Get().(*Event)
	e.closed = false
	e.id = gonanoid.Must(16)
	e.createdAt = time.Now()
	e.emittedAt = time.Time{}
	e.level = level
	e.message = msg
	e.topic = ""
	e.parentID = ""
	e.fields = make(Fields)
	e.c = c
	return e
}

func newChildEvent(level Level, msg string, e *Event) *Event {
	if e == nil || e.c == nil || level < e.c.level {
		return nil
	}
	ce := newEvent(level, msg, e.c)
	return ce.ParentID(e.id)
}

// Trace creates a child event with the trace level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new child event.
func (e *Event) Trace(args ...any) *Event {
	return newChildEvent(LevelTrace, fmt.Sprint(args...), e)
}

// Tracef creates a child event with the trace level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new child event.
func (e *Event) Tracef(v string, args ...any) *Event {
	return newChildEvent(LevelTrace, fmt.Sprintf(v, args...), e)
}

// Debug creates a child event with the debug level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new child event.
func (e *Event) Debug(args ...any) *Event {
	return newChildEvent(LevelDebug, fmt.Sprint(args...), e)
}

// Debugf creates a child event with the debug level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new child event.
func (e *Event) Debugf(v string, args ...any) *Event {
	return newChildEvent(LevelDebug, fmt.Sprintf(v, args...), e)
}

// Info creates a child event with the info level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new child event.
func (e *Event) Info(args ...any) *Event {
	return newChildEvent(LevelInfo, fmt.Sprint(args...), e)
}

// Infof creates a child event with the info level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new child event.
func (e *Event) Infof(v string, args ...any) *Event {
	return newChildEvent(LevelInfo, fmt.Sprintf(v, args...), e)
}

// Warn creates a child event with the warn level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new child event.
func (e *Event) Warn(args ...any) *Event {
	return newChildEvent(LevelWarn, fmt.Sprint(args...), e)
}

// Warnf creates a child event with the warn level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new child event.
func (e *Event) Warnf(v string, args ...any) *Event {
	return newChildEvent(LevelWarn, fmt.Sprintf(v, args...), e)
}

// Error creates a child event with the error level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new child event.
func (e *Event) Error(args ...any) *Event {
	return newChildEvent(LevelError, fmt.Sprint(args...), e)
}

// Errorf creates a child event with the error level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new child event.
func (e *Event) Errorf(v string, args ...any) *Event {
	return newChildEvent(LevelError, fmt.Sprintf(v, args...), e)
}

// Fatal creates a child event with the fatal level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new child event.
func (e *Event) Fatal(args ...any) *Event {
	return newChildEvent(LevelFatal, fmt.Sprint(args...), e)
}

// Fatalf creates a child event with the fatal level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new child event.
func (e *Event) Fatalf(v string, args ...any) *Event {
	return newChildEvent(LevelFatal, fmt.Sprintf(v, args...), e)
}

// Panic creates a child event with the panic level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprint and assigns this message to the new child event.
func (e *Event) Panic(args ...any) *Event {
	return newChildEvent(LevelPanic, fmt.Sprint(args...), e)
}

// Panicf creates a child event with the panic level and sets the parentID to the current event's ID.
// It formats the provided arguments into a message using fmt.Sprintf according to the format string 'v' and assigns this message to the new child event.
func (e *Event) Panicf(v string, args ...any) *Event {
	return newChildEvent(LevelPanic, fmt.Sprintf(v, args...), e)
}

func (e *Event) String() string {
	return fmt.Sprintf(
		"&{ID:%s CreatedAt:%v Level:%s Message:%s Topic:%s ParentID:%s Fields:%v}",
		e.id,
		e.createdAt,
		e.level,
		e.message,
		e.topic,
		e.parentID,
		e.fields,
	)
}

func (e *Event) ID() string {
	if e == nil {
		return ""
	}
	return e.id
}

// E is an alias for the Emit method.
// Emit triggers the event, notifying all observers that match the condition.
// If hold is true, the event is not recycled and must be manually closed with `Evict`. Otherwise, the event is returned to the pool after emission.
func (e *Event) E(hold ...bool) *Event { return e.Emit(hold...) }

// Emit triggers the event, notifying all observers that match the condition.
// If hold is true, the event is not recycled and must be manually closed with `Evict`. Otherwise, the event is returned to the pool after emission.
func (e *Event) Emit(hold ...bool) *Event {
	if e == nil {
		return nil
	}
	if e.closed {
		return e
	}

	e.emittedAt = time.Now()
	e.closed = true

	for _, o := range e.c.observers {
		if o.cond(e) {
			o.handler(e)
		}
	}

	if len(hold) > 0 && hold[0] {
		return e
	}

	eventPool.Put(e)
	return e
}

func (e *Event) Evict() {
	if e == nil {
		return
	}
	eventPool.Put(e)
}

func (e *Event) WithError(err error) *Event {
	if e == nil {
		return e
	}
	return e.Fields(Fields{"error": err})
}

// F is an alias for the Field method.
func (e *Event) F(k string, v any) *Event { return e.Field(k, v) }

func (e *Event) Field(k string, v any) *Event {
	if e == nil {
		return e
	}
	return e.Fields(Fields{k: v})
}

// Fs is an alias for the Fields method.
func (e *Event) Fs(fields Fields) *Event { return e.Fields(fields) }

func (e *Event) Fields(fields Fields) *Event {
	if e == nil {
		return e
	}
	data := make(Fields)
	for k, v := range e.fields {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	e.fields = data
	return e
}

func (e *Event) Level(level Level) *Event {
	if e == nil {
		return e
	}
	e.level = level
	return e
}

// M is an alias for the Message method.
func (e *Event) M(args ...any) *Event { return e.Message(args...) }

func (e *Event) Message(args ...any) *Event {
	if e == nil {
		return e
	}
	e.message = fmt.Sprint(args...)
	return e
}

// Mf is an alias for the Messagef method.
func (e *Event) Mf(f string, args ...any) *Event { return e.Messagef(f, args...) }

func (e *Event) Messagef(f string, args ...any) *Event {
	if e == nil {
		return e
	}
	e.message = fmt.Sprintf(f, args...)
	return e
}

// P is an alias for the ParentID method.
func (e *Event) P(id string) *Event { return e.ParentID(id) }

func (e *Event) ParentID(id string) *Event {
	if e == nil {
		return e
	}
	e.parentID = id
	return e
}

// T is an alias for the Topic method.
func (e *Event) T(topic string) *Event { return e.Topic(topic) }

func (e *Event) Topic(topic string) *Event {
	if e == nil {
		return e
	}
	e.topic = topic
	return e
}
