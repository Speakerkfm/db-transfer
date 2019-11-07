package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
)

var (
	consoleBufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
)

func ReportPortalWriter(endpoint, project, token string) io.Writer {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00"

	endpoint += "/" + project + "/log"
	return rpWriter{endpoint: endpoint, token: token}
}

type rpMessage struct {
	ItemID  interface{} `json:"item_id"`
	Level   interface{} `json:"level"`
	Message string      `json:"message"`
	Time    interface{} `json:"time"`
}

type rpWriter struct {
	endpoint string
	token    string
}

func (w rpWriter) Write(p []byte) (n int, err error) {
	//spew.Dump(string(p))

	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	var evt map[string]interface{}

	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	reqID, ok := evt[RequestIDFieldName]
	if !ok {
		return
	}
	level, ok := evt[zerolog.LevelFieldName]
	if !ok {
		return
	}
	if level == "panic" || level == "fatal" {
		level = "error"
	}
	tt, ok := evt[zerolog.TimestampFieldName]
	if !ok {
		return
	}

	msg, _ := json.MarshalIndent(evt, "", "    ")
	err = json.NewEncoder(buf).Encode(&rpMessage{
		ItemID:  reqID,
		Level:   level,
		Message: string(msg),
		Time:    tt,
	})

	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", w.endpoint, buf)
	req.Header.Add("Authorization", "bearer "+w.token)
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)

	//b, _ := httputil.DumpRequest(req, true)
	//b, _ := httputil.DumpResponse(resp, true)
	//spew.Dump(evt)
	//spew.Dump(string(b))
	//spew.Dump(string(w.endpoint))
	//spew.Dump(buf)
	return len(p), err
}
