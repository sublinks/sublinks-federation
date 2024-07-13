package log

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
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

func NewLogger(name string) *Log {
	logger := zerolog.New(os.Stdout)
	logger.Debug().Msg(fmt.Sprintf("%s logger started", name))
	return &Log{&logger}
}

func SetGlobalLevel(level string) {
	switch strings.ToLower(level) {
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func (logger *Log) Info(msg string) {
	logger.Logger.Info().Msg(msg)
}

func (logger *Log) Debug(msg string) {
	logger.Logger.Debug().Msg(msg)
}

func (logger *Log) Error(msg string, err error) {
	logger.Logger.Error().Err(err).Msg(msg)
}

func (logger *Log) Fatal(msg string, err error) {
	logger.Logger.Fatal().Err(err).Msg(msg)
}

func (logger *Log) Warn(msg string) {
	logger.Logger.Warn().Msg(msg)
}

func (logger *Log) Request(msg string, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("error closing http io reader")
		}
	}(r.Body)
	var body interface{}
	rawbody, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Error reading request body", err)
		body = nil
	}
	if r.ContentLength > 0 && r.Header.Get("Content-Type") == "application/json" {
		err = json.Unmarshal(rawbody, &body)
		if err != nil {
			logger.Error("Error parsing request body into json", err)
			body = nil
		}
	} else {
		body = rawbody
	}
	logger.Logger.Debug().
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
