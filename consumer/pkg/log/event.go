package log

import (
	"github.com/rs/zerolog"
)

type Event struct {
	logEvent  *zerolog.Event
	dictEvent *zerolog.Event
}

func (e *Event) Msg(val string) {
	if e == nil {
		return
	}
	if e.logEvent == nil {
		return
	}
	if e.dictEvent != nil {
		e.logEvent.Dict("context", e.dictEvent)
	}
	e.logEvent.Msg(val)
}

func (e *Event) User(val string) *Event {
	if e == nil {
		return e
	}
	e.logEvent.Str("LOGIN_USER_ID", val)
	return e
}

func (e *Event) Project(val string) *Event {
	if e == nil {
		return e
	}
	e.logEvent.Str("LOGIN_PROJECT_ID", val)
	return e
}
