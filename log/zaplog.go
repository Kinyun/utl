package log

import (
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"sync"
	"time"
)

var instance *zap.Logger
var sugarInstance *zap.SugaredLogger
var once sync.Once

func getLogger() *zap.Logger {
	once.Do(func() {
		instance = newLogger()
	})
	return instance
}

func GetSugaredLogger() *zap.SugaredLogger {
	if sugarInstance == nil {
		sugarInstance = getLogger().Sugar()
	}
	return sugarInstance
}

func newLogger() *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout"}
	if logger, err := loggerConfig.Build(); err != nil {
		Errorf("failed to create new logger with error: %s", err)
		panic(err)
	} else {
		return logger
	}
}

// Debug logs the message at debug level with additional fields, if any
func Debug(message string, fields ...zap.Field) {
	getLogger().Debug(message, fields...)
}

// Debugf allows Sprintf style formatting and logs at debug level
func Debugf(template string, args ...interface{}) {
	GetSugaredLogger().Debugf(template, args...)
}

// Error logs is equivalent to Error() followed by a call to panic().
func Panic(message string, fields ...zap.Field) {
	getLogger().Error(message, fields...)
	panic(message)
}

// Error logs the message at error level and prints stacktrace with additional fields, if any
func Error(err error, fields ...zap.Field) {
	getLogger().Error(err.Error(), fields...)
}

// Error logs the message at error level and prints stacktrace with additional fields, if any
func SError(err string, fields ...zap.Field) {
	getLogger().Error(err, fields...)
}

// Errorf allows Sprintf style formatting, logs at error level and prints stacktrace
func Errorf(template string, args ...interface{}) {
	GetSugaredLogger().Errorf(template, args...)
}

// Fatal logs the message at fatal level with additional fields, if any and exits
func Fatal(err error, fields ...zap.Field) {
	getLogger().Fatal(err.Error(), fields...)
}

// Fatalf allows Sprintf style formatting, logs at fatal level and exits
func Fatalf(template string, args ...interface{}) {
	GetSugaredLogger().Fatalf(template, args...)
}

// Info logs the message at info level with additional fields, if any
func Info(message string, fields ...zap.Field) {
	getLogger().Info(message, fields...)
}

// Infof allows Sprintf style formatting and logs at info level
func Infof(template string, args ...interface{}) {
	GetSugaredLogger().Infof(template, args...)
}

// Warn logs the message at warn level with additional fields, if any
func Warn(message string, fields ...zap.Field) {
	getLogger().Warn(message, fields...)
}

// Warnf allows Sprintf style formatting and logs at warn level
func Warnf(template string, args ...interface{}) {
	GetSugaredLogger().Warnf(template, args...)
}

// AddHook adds func(zapcore.Entry) error) to the logger lifecycle
func AddHook(hook func(zapcore.Entry) error) {
	instance = getLogger().WithOptions(zap.Hooks(hook))
	sugarInstance = instance.Sugar()
}

// WithRequest takes in a http.Request and logs the message with request's Method, Host and Path
// and returns zap.logger
func WithRequest(r *http.Request) *zap.Logger {
	return getLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}

// SugaredWithRequest takes in a http.Request and logs the message with request's Method, Host and Path
// and returns zap.SugaredLogger to support Sprintf styled logging
func SugaredWithRequest(r *http.Request) *zap.SugaredLogger {
	return GetSugaredLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}

func LatencyLogger(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(negroni.ResponseWriter)

	getLogger().Info(
		"LatencyLogger",
		zap.Int("status", res.Status()),
		zap.String("start", start.Format(time.RFC1123Z)),
		zap.String("latency", time.Since(start).String()),
		zap.String("hostname", r.Host),
		zap.String("user_agent", r.UserAgent()),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("parameter", r.URL.RawQuery),
	)
}
