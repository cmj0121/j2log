package j2log

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

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
	line = fmt.Sprintf("%s %s", timestamp, message)
	ok = true
	return
}

func (t Template) extract(raw map[string]interface{}, key string) (value string, ok bool) {
	if value, ok = raw[key].(string); !ok {
		log.Trace().Str("key", key).Msg("failed to extract")
		return
	}

	return
}
