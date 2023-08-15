package j2log

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
)

// the extracted log
type Log struct {
	Timestamp time.Time
	Message   string

	// the rest of the log
	Rest map[string]interface{}
}

// the template of the final human readable log.
type Template struct {
	Timestamp string
	Message   string
}

// create the default template
func DefaultTmpl() *Template {
	return &Template{
		Timestamp: "@timestamp",
		Message:   "message",
	}
}

// convert the JSON log to human-readable log
func (t Template) Extract(raw map[string]interface{}) (line string, ok bool) {
	var extracted Log

	if extracted, ok = t.extractLog(raw); !ok {
		log.Info().Msg("failed to extract log")
		return
	}

	// construct the human-readable log
	var buff bytes.Buffer

	tmpl := template.Must(template.New("log").Parse(`[{{- .Timestamp.Format "2006-01-02T15:04:05-0700" -}}] {{ .Message -}}`))
	if err := tmpl.Execute(&buff, extracted); err != nil {
		log.Debug().Err(err).Msg("failed to execute template")
		return
	}

	line = buff.String()
	ok = true
	return
}

func (t Template) extractLog(raw map[string]interface{}) (extracted Log, ok bool) {
	var timestamp, message string

	if timestamp, ok = t.extract(raw, t.Timestamp); !ok {
		log.Debug().Msg("failed to extract timestamp")
		return
	}

	if message, ok = t.extract(raw, t.Message); !ok {
		log.Debug().Msg("failed to extract message")
		return
	}

	// construct the human-readable log
	var layout string
	switch {
		case strings.HasSuffix(timestamp, "Z"):
			layout = "2006-01-02T15:04:05.999Z"
		default:
			layout = "2006-01-02T15:04:05.999-0700"
	}

	tz, err := time.Parse(layout, timestamp)
	if err != nil {
		log.Debug().Err(err).Msg("failed to parse timestamp")
		return
	}

	delete(raw, t.Timestamp)
	delete(raw, t.Message)

	extracted = Log{
		Timestamp: tz,
		Message:   message,
		Rest:      raw,
	}
	return
}

func (t Template) extract(raw map[string]interface{}, key string) (value string, ok bool) {
	if value, ok = raw[key].(string); !ok {
		log.Trace().Str("key", key).Msg("failed to extract")
		return
	}

	return
}
