package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	grpc_zerolog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
)

const (
	miliseconds float32 = 1000000.0
)

var (
	once sync.Once
	log  *Log
)

func GetLogger() *Log {
	once.Do(func() {
		log = NewLog()
	})
	return log
}

type Log struct {
	logger zerolog.Logger
}

func NewLog() *Log {
	return NewLogWithLevel(getLogLevel("LOG_LEVEL"))
}

func NewLogWithLevel(level string) *Log {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	return &Log{logger: zerolog.New(os.Stdout).Level(lvl).With().Timestamp().Logger()}

}

func (l *Log) GetLogger() grpc_zerolog.Logger {
	return grpc_zerolog.LoggerFunc(func(ctx context.Context, level grpc_zerolog.Level, msg string, fields ...any) {
		l := l.logger.With().Fields(fields).Logger()

		switch level {
		case grpc_zerolog.LevelDebug:
			l.Debug().Msg(msg)
		case grpc_zerolog.LevelInfo:
			l.Info().Msg(msg)
		case grpc_zerolog.LevelWarn:
			l.Warn().Msg(msg)
		case grpc_zerolog.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", level))
		}
	})
}

func (l *Log) GetOptions() []grpc_zerolog.Option {
	return []grpc_zerolog.Option{
		grpc_zerolog.WithLevels(codeToLevel),
		grpc_zerolog.WithLogOnEvents(grpc_zerolog.FinishCall),
	}
}

func (l *Log) Debug(moduleName, functionName, msg string) {
	l.logger.Debug().
		Str("module", moduleName).
		Str("function", functionName).
		Msg(msg)
}

func (l *Log) Debugf(format string, args ...any) {
	l.logger.Debug().Msgf(format, args...)
}

func (l *Log) Infof(format string, args ...any) {
	l.logger.Info().Msgf(format, args...)
}

func (l *Log) Warning(moduleName, functionName, msg string) {
	l.logger.Warn().
		Str("module", moduleName).
		Str("function", functionName).
		Msg(msg)
}

func (l *Log) Warningf(moduleName, functionName, format string, args ...any) {
	l.logger.Warn().
		Str("module", moduleName).
		Str("function", functionName).
		Msgf(format, args...)
}

func (l *Log) Error(moduleName, functionName string, err error) {
	l.logger.Error().
		Str("module", moduleName).
		Str("function", functionName).
		Msg(err.Error())
}

func (l *Log) Errorf(moduleName, functionName, format string, args ...any) {
	l.logger.Error().
		Str("module", moduleName).
		Str("function", functionName).
		Msgf(format, args...)
}

func (l *Log) FatalIfError(moduleName, functionName string, errs ...error) {
	var sb strings.Builder
	for _, err := range errs {
		if err != nil {
			sb.WriteString(err.Error())
			sb.WriteString("\t")
		}
	}
	if sb.Len() > 0 {
		l.Fatalf(moduleName, functionName, sb.String())
	}
}

func (l *Log) Fatal(moduleName, functionName string, err error) {
	l.Fatalf(moduleName, functionName, err.Error())
}

func (l *Log) Fatalf(moduleName, functionName, format string, args ...any) {
	l.logger.Fatal().
		Str("module", moduleName).
		Str("function", functionName).
		Msgf(format, args...)
}

func (l *Log) Request(context echo.Context, start time.Time) {
	if context.Get("response-error") == nil {
		l.logger.Info().
			Str("method", context.Request().Method).
			Int("status", context.Response().Status).
			Str("request", context.Request().RequestURI).
			Interface("responseBody", context.Get("response-body")).
			Float32("timestamp", timestamp(start)).
			Msg("")
	} else {
		l.logger.Error().
			Str("method", context.Request().Method).
			Int("status", context.Response().Status).
			Str("request", context.Request().RequestURI).
			Interface("responseBody", context.Get("response-body")).
			Str("error", context.Get("response-error").(string)).
			Float32("timestamp", timestamp(start)).
			Msg("")
	}
}

func timestamp(start time.Time) float32 {
	return float32(time.Since(start).Nanoseconds()) / miliseconds
}

func getLogLevel(logLevel string) string {
	value, found := os.LookupEnv(logLevel)
	if !found {
		return "debug"
	}
	return value
}

func codeToLevel(code codes.Code) grpc_zerolog.Level {
	if code == codes.OK {
		return grpc_zerolog.LevelDebug
	}
	return grpc_zerolog.DefaultClientCodeToLevel(code)
}
