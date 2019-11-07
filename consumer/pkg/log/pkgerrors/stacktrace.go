package pkgerrors

import (
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

var (
	StackSourceFileName     = "source"
	StackSourceLineName     = "line"
	StackSourceFunctionName = "func"
	StackSourceLen          = 7
)

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b
	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(c int) bool {
	return true
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)
	return string(s.b)
}

// MarshalStack implements pkg/errors stack trace marshaling.
//
//   zerolog.ErrorStackMarshaler = MarshalStack
func MarshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	sterr, ok := err.(stackTracer)
	if !ok {
		return parseStack(string(Stack(4, StackSourceLen)))
	}
	st := sterr.StackTrace()
	s := &state{}
	out := make([]map[string]string, 0, len(st))

	stackLen := len(st)
	if stackLen > StackSourceLen {
		stackLen = StackSourceLen
	}

	for _, frame := range st[0:stackLen] {
		filePath := frameField(frame, s, 's')
		if sepIndex := strings.Index(filePath, "\n\t"); sepIndex > 0 {
			filePath = filePath[sepIndex+2:]
		}

		out = append(out, map[string]string{
			StackSourceFileName:     filePath,
			StackSourceLineName:     frameField(frame, s, 'd'),
			StackSourceFunctionName: frameField(frame, s, 'n'),
		})
	}

	return out
}

func parseStack(text string) []map[string]string {
	out := make([]map[string]string, 0, len(text)/2)

	lines := strings.Split(text, "\n")

	var (
		frame      map[string]string
		colonIndex int
		spaceIndex int
	)
	for i, l := range lines {
		if i%2 == 0 {
			colonIndex = strings.LastIndex(l, "(")
			if colonIndex <= 0 {
				colonIndex = len(l)
			}
			frame = map[string]string{
				StackSourceFunctionName: l[0:colonIndex],
				StackSourceFileName:     "",
				StackSourceLineName:     "",
			}
		} else {
			colonIndex = strings.LastIndex(l, ":")
			spaceIndex = strings.LastIndex(l, " ")
			if spaceIndex == -1 {
				spaceIndex = len(l)
			}
			frame[StackSourceFileName] = strings.TrimPrefix(l[0:colonIndex], "\t")
			frame[StackSourceLineName] = l[colonIndex+1 : spaceIndex]
			out = append(out, frame)
		}
	}

	return out
}

func Stack(skipCallers, callersNum int) []byte {
	var (
		e             = make([]byte, 1<<16) // 64k
		nbytes        = runtime.Stack(e, false)
		ignorelinenum = 2*skipCallers + 1
		count         = 0
		startIndex    = 0
		stopIndex     = 0
		traceLinesNum = 2*callersNum + ignorelinenum
	)

	for i := 0; i <= nbytes; i++ {
		if e[i] == '\n' {
			count++
			if count == ignorelinenum {
				startIndex = i + 1
			}

			if startIndex > 0 && count == traceLinesNum {
				stopIndex = i + 1
			}
		}
	}

	if stopIndex == 0 {
		stopIndex = nbytes
	}

	return e[startIndex:stopIndex]
}
