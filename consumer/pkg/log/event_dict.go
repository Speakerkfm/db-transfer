package log

import (
	"net"
	"time"

	"github.com/rs/zerolog"
)

func (e *Event) checkDictEventExist() {
	if e.dictEvent == nil {
		e.dictEvent = zerolog.Dict()
	}
}

func (e *Event) Err(err error) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Stack().Err(err)
	return e
}

// Errs adds the field key with errs as an array of serialized errors to the
// *Event context.
func (e *Event) Errs(key string, errs []error) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Stack().Errs(key, errs)
	return e
}

// Fields is a helper function to use a map to set fields using type assertion.
func (e *Event) Fields(fields map[string]interface{}) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Fields(fields)
	return e
}

// Array adds the field key with an array to the event context.
// Use zerolog.Arr() to create the array or pass a type that
// implement the LogArrayMarshaler interface.
func (e *Event) Array(key string, arr zerolog.LogArrayMarshaler) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Array(key, arr)
	return e
}

// Object marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) Object(key string, obj zerolog.LogObjectMarshaler) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Object(key, obj)
	return e
}

// Object marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) EmbedObject(obj zerolog.LogObjectMarshaler) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.EmbedObject(obj)
	return e
}

// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Str(key, val)
	return e
}

// Strs adds the field key with vals as a []string to the *Event context.
func (e *Event) Strs(key string, vals []string) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Strs(key, vals)
	return e
}

// Bytes adds the field key with val as a string to the *Event context.
//
// Runes outside of normal ASCII ranges will be hex-encoded in the resulting
// JSON.
func (e *Event) Bytes(key string, val []byte) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Bytes(key, val)
	return e
}

// Hex adds the field key with val as a hex string to the *Event context.
func (e *Event) Hex(key string, val []byte) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Hex(key, val)
	return e
}

// RawJSON adds already encoded JSON to the log line under key.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (e *Event) RawJSON(key string, b []byte) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.RawJSON(key, b)
	return e
}

// Bool adds the field key with val as a bool to the *Event context.
func (e *Event) Bool(key string, b bool) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Bool(key, b)
	return e
}

// Bools adds the field key with val as a []bool to the *Event context.
func (e *Event) Bools(key string, b []bool) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Bools(key, b)
	return e
}

// Int adds the field key with i as a int to the *Event context.
func (e *Event) Int(key string, i int) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Int(key, i)
	return e
}

// Ints adds the field key with i as a []int to the *Event context.
func (e *Event) Ints(key string, i []int) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Ints(key, i)
	return e
}

// Int8 adds the field key with i as a int8 to the *Event context.
func (e *Event) Int8(key string, i int8) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Int8(key, i)
	return e
}

// Ints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Ints8(key string, i []int8) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Ints8(key, i)
	return e
}

// Int16 adds the field key with i as a int16 to the *Event context.
func (e *Event) Int16(key string, i int16) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Int16(key, i)
	return e
}

// Ints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Ints16(key string, i []int16) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Ints16(key, i)
	return e
}

// Int32 adds the field key with i as a int32 to the *Event context.
func (e *Event) Int32(key string, i int32) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Int32(key, i)
	return e
}

// Ints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Ints32(key string, i []int32) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Ints32(key, i)
	return e
}

// Int64 adds the field key with i as a int64 to the *Event context.
func (e *Event) Int64(key string, i int64) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Int64(key, i)
	return e
}

// Ints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Ints64(key string, i []int64) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Ints64(key, i)
	return e
}

// Uint adds the field key with i as a uint to the *Event context.
func (e *Event) Uint(key string, i uint) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uint(key, i)
	return e
}

// Uints adds the field key with i as a []int to the *Event context.
func (e *Event) Uints(key string, i []uint) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uints(key, i)
	return e
}

// Uint8 adds the field key with i as a uint8 to the *Event context.
func (e *Event) Uint8(key string, i uint8) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uint8(key, i)
	return e
}

// Uints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Uints8(key string, i []uint8) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uints8(key, i)
	return e
}

// Uint16 adds the field key with i as a uint16 to the *Event context.
func (e *Event) Uint16(key string, i uint16) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uint16(key, i)
	return e
}

// Uints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Uints16(key string, i []uint16) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uints16(key, i)
	return e
}

// Uint32 adds the field key with i as a uint32 to the *Event context.
func (e *Event) Uint32(key string, i uint32) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uint32(key, i)
	return e
}

// Uints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Uints32(key string, i []uint32) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uints32(key, i)
	return e
}

// Uint64 adds the field key with i as a uint64 to the *Event context.
func (e *Event) Uint64(key string, i uint64) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uint64(key, i)
	return e
}

// Uints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Uints64(key string, i []uint64) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Uints64(key, i)
	return e
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e *Event) Float32(key string, f float32) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Float32(key, f)
	return e
}

// Floats32 adds the field key with f as a []float32 to the *Event context.
func (e *Event) Floats32(key string, f []float32) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Floats32(key, f)
	return e
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e *Event) Float64(key string, f float64) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Float64(key, f)
	return e
}

// Floats64 adds the field key with f as a []float64 to the *Event context.
func (e *Event) Floats64(key string, f []float64) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Floats64(key, f)
	return e
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *Event) Time(key string, t time.Time) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Time(key, t)
	return e
}

// Times adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *Event) Times(key string, t []time.Time) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Times(key, t)
	return e
}

// Dur adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Dur(key string, d time.Duration) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Dur(key, d)
	return e
}

// Durs adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Durs(key string, d []time.Duration) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Durs(key, d)
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.TimeDiff(key, t, start)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *Event) Interface(key string, i interface{}) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.Interface(key, i)
	return e
}

// Caller adds the file:line of the caller with the zerolog.CallerFieldName key.
//func (e *Event) Caller() *Event {
//	return e.caller(CallerSkipFrameCount)
//}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *Event) IPAddr(key string, ip net.IP) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.IPAddr(key, ip)
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.IPPrefix(key, pfx)
	return e
}

// MACAddr adds MAC address to the event
func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	if e == nil {
		return e
	}
	e.checkDictEventExist()
	e.dictEvent.MACAddr(key, ha)
	return e
}
