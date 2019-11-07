package log

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	RequestIDFieldName = "uid"
	programFieldName,
	programFieldValue string
)

func Config(programName string, writer io.Writer) {
	programFieldValue = programName
	programFieldName = "PROGRAM"
	zerolog.MessageFieldName = "MESSAGE"
	zerolog.LevelFieldName = "LEVEL"
	log.Logger = log.Output(writer)
}

type Logger struct {
	r             *http.Request
	zerologLogger *zerolog.Logger
}

func FromRequest(r *http.Request) *Logger {
	return NewLogger(r)
}

func (l *Logger) Debug() *Event {
	return withLevel(l, zerolog.DebugLevel)
}

func (l *Logger) Info() *Event {
	return withLevel(l, zerolog.InfoLevel)
}

func (l *Logger) Warn() *Event {
	return withLevel(l, zerolog.WarnLevel)
}

func (l *Logger) Error() *Event {
	return withLevel(l, zerolog.ErrorLevel)
}

func (l *Logger) Fatal() *Event {
	return withLevel(l, zerolog.FatalLevel)
}

func (l *Logger) Panic() *Event {
	return withLevel(l, zerolog.PanicLevel)
}

func Debug() *Event {
	l := NewLogger(nil)
	return withLevel(l, zerolog.DebugLevel)
}

func Info() *Event {
	l := NewLogger(nil)
	return withLevel(l, zerolog.InfoLevel)
}

func Warn() *Event {
	l := NewLogger(nil)
	return withLevel(l, zerolog.WarnLevel)
}

func Error() *Event {
	l := NewLogger(nil)
	return withLevel(l, zerolog.ErrorLevel)
}

func Fatal() *Event {
	l := NewLogger(nil)
	return withLevel(l, zerolog.FatalLevel)
}

func Panic() *Event {
	l := NewLogger(nil)
	return withLevel(l, zerolog.PanicLevel)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Debug().Msg(fmt.Sprintf(format, v...))
}

func withLevel(l *Logger, level zerolog.Level) *Event {
	event := l.zerologLogger.WithLevel(level)
	return &Event{logEvent: event.Str(programFieldName, programFieldValue)}
}

func NewLogger(r *http.Request) *Logger {
	var zerologLogger *zerolog.Logger

	if r != nil {
		zerologLogger = log.Ctx(r.Context())
	} else {
		zerologLogger = &log.Logger
	}

	return &Logger{r: r, zerologLogger: zerologLogger}
}

func Dict() *zerolog.Event {
	return zerolog.Dict()
}
