package log

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func init() {
	log.Debug().Msg("Logger started")
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Error(msg string, err error) {
	log.Error().Err(err).Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Request(msg string, r *http.Request) {
	defer r.Body.Close()
	var body interface{}
	rawbody, err := io.ReadAll(r.Body)
	if err != nil {
		Error("Error reading request body", err)
		body = nil
	}
	if r.ContentLength > 0 && r.Header.Get("Content-Type") == "application/json" {
		err = json.Unmarshal(rawbody, &body)
		if err != nil {
			Error("Error parsing request body into json", err)
			body = nil
		}
	}
	log.Debug().
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Str("user-agent", r.UserAgent()).
		Int64("content-length", r.ContentLength).
		Str("ip", r.RemoteAddr).
		Str("real-ip", r.Header.Get("X-Real-Ip")).
		Str("dnt", r.Header.Get("Dnt")).
		Str("host", r.Host).
		Str("proto", r.Proto).
		Str("referer", r.Referer()).
		Str("accept", r.Header.Get("Accept")).
		Str("accept-language", r.Header.Get("Accept-Language")).
		Str("content-type", r.Header.Get("Content-Type")).
		Str("forwarded-for", r.Header.Get("X-Forwarded-For")).
		Any("body", body).
		Str("user", r.URL.User.Username()).
		Str("query", r.URL.Query().Encode()).
		Msg(msg)
}
