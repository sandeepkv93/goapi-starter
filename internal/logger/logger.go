package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	requestLogger    zerolog.Logger
	requestLoggerMux sync.RWMutex
)

func Init() {
	// Pretty print during development
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Set caller marshaler to short path
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		return short
	}

	log.Logger = log.With().Caller().Logger()
}

// SetRequestLogger sets the logger for the current request
func SetRequestLogger(l zerolog.Logger) {
	requestLoggerMux.Lock()
	defer requestLoggerMux.Unlock()
	requestLogger = l
}

// ClearRequestLogger clears the request logger
func ClearRequestLogger() {
	requestLoggerMux.Lock()
	defer requestLoggerMux.Unlock()
	requestLogger = log.Logger // Reset to default logger
}

// Debug returns a debug level event logger with request context
func Debug() *zerolog.Event {
	requestLoggerMux.RLock()
	defer requestLoggerMux.RUnlock()
	return requestLogger.Debug()
}

// Info returns an info level event logger with request context
func Info() *zerolog.Event {
	requestLoggerMux.RLock()
	defer requestLoggerMux.RUnlock()
	return requestLogger.Info()
}

// Warn returns a warn level event logger with request context
func Warn() *zerolog.Event {
	requestLoggerMux.RLock()
	defer requestLoggerMux.RUnlock()
	return requestLogger.Warn()
}

// Error returns an error level event logger with request context
func Error() *zerolog.Event {
	requestLoggerMux.RLock()
	defer requestLoggerMux.RUnlock()
	return requestLogger.Error()
}

// Fatal returns a fatal level event logger with request context
func Fatal() *zerolog.Event {
	requestLoggerMux.RLock()
	defer requestLoggerMux.RUnlock()
	return requestLogger.Fatal()
}
