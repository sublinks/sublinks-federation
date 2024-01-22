package log

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string, err error)
	Fatal(msg string, err error)
	Warn(msg string)
	Request(msg string, r *http.Request)
}

type Log struct {
	*zerolog.Logger
}

func NewLogger() *Log {
	log.Debug().Msg("Logger started")
	return &Log{}
}

func (l Log) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l Log) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}

func (l Log) Error(msg string, err error) {
	l.Logger.Error().Err(err).Msg(msg)
}

func (l Log) Fatal(msg string, err error) {
	l.Logger.Fatal().Err(err).Msg(msg)
}

func (l Log) Warn(msg string) {
	l.Logger.Warn().Msg(msg)
}

func (l Log) Request(msg string, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("error closing http io reader")
		}
	}(r.Body)
	var body interface{}
	rawbody, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("Error reading request body", err)
		body = nil
	}
	if r.ContentLength > 0 && r.Header.Get("Content-Type") == "application/json" {
		err = json.Unmarshal(rawbody, &body)
		if err != nil {
			l.Error("Error parsing request body into json", err)
			body = nil
		}
	} else {
		body = rawbody
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
