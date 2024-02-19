package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger interface {
		Debug(template string, args ...interface{})
		Info(template string, args ...interface{})
		Warn(template string, args ...interface{})
		Error(template string, args ...interface{})

		// Panic with message
		// If error is wrapped with errors.Wrap method, use format %+v to print stack trace
		Panic(format string, args ...interface{})

		// Generic log printer with input level manually
		Print(level string, fields map[string]interface{}, template string, args ...interface{})
		WithModule(name string) Logger
		ContextLogger
	}

	Config struct {
		// Log level.
		// Can be one of: debug, info, warn, error, panic
		Level string `mapstructure:"level"`
		// Where log will be written to.
		// Can be one of: stdout, stderr, file path
		Output []string `mapstructure:"output"`

		// Where to write error output: eg. stderr
		ErrOutput []string `mapstructure:"errOutput"`

		// Log output format.
		// Can be one of: console, json
		Format string `mapstructure:"format"`
	}

	LogLevel = string

	entry struct {
		inst *zap.SugaredLogger
	}

	LogFieldsFn = func(msg string, fields ...interface{})
)

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
)

var inst *zap.Logger

// Initiate logger instance with default Config
func init() {
	InitLogInst(Config{
		Level:     "debug",
		Output:    []string{"stdout"},
		ErrOutput: []string{"stderr"},
		Format:    "json",
	})
	inst.Debug("initiate default log engine")
}

// Define global log inst using input Config
func InitLogInst(conf Config) {
	level := zap.InfoLevel
	if err := level.Set(conf.Level); err != nil {
		panic(errors.WithMessage(err, "invalid log level config"))
	}
	encoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "module",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	zapConf := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		DisableCaller:    true,
		Sampling:         nil,
		Encoding:         conf.Format,
		EncoderConfig:    encoder,
		OutputPaths:      conf.Output,
		ErrorOutputPaths: conf.ErrOutput,
	}
	log, err := zapConf.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(errors.WithMessage(err, "failed to initiate log instance"))
	}
	log.Info(fmt.Sprintf("log engine load successfully with config: %+v", conf))
	inst = log
}

func (e *entry) Print(level LogLevel, fields map[string]interface{}, template string, args ...interface{}) {
	var fn LogFieldsFn
	switch level {
	case InfoLevel:
		fn = e.inst.Infow
	case WarnLevel:
		fn = e.inst.Warnw
	case ErrorLevel:
		fn = e.inst.Errorw
	case PanicLevel:
		fn = e.inst.Panicw
	default:
		fn = e.inst.Debugw
	}
	var keyAndValues []interface{}
	for k, v := range fields {
		keyAndValues = append(keyAndValues, k, v)
	}
	fn(fmt.Sprintf(template, args...), keyAndValues...)
}

func (e *entry) Debug(template string, args ...interface{}) {
	e.inst.Debugf(template, args...)
}

func (e *entry) Info(template string, args ...interface{}) {
	e.inst.Infof(template, args...)
}

func (e *entry) Warn(template string, args ...interface{}) {
	e.inst.Warnf(template, args...)
}

func (e *entry) Error(template string, args ...interface{}) {
	e.inst.Errorf(template, args...)
}

func (e *entry) Panic(template string, args ...interface{}) {
	e.inst.Panicf(template, args...)
}

func (e *entry) Printf(format string, args ...interface{}) {
	e.inst.Debugf(format, args...)
}

func (e *entry) WithModule(name string) Logger {
	return &entry{inst: e.inst.Named(name)}
}

// Initiate a log instance from the global core logger
func WithModule(name string) Logger {
	logger := inst.Sugar().Named(name)
	return &entry{inst: logger}
}

func Info(template string, args ...interface{}) {
	inst.Info(fmt.Sprintf(template, args...))
}

func Panic(template string, args ...interface{}) {
	inst.Panic(fmt.Sprintf(template, args...))
}

func Error(template string, args ...interface{}) {
	inst.Error(fmt.Sprintf(template, args...))
}

func Debug(template string, args ...interface{}) {
	inst.Debug(fmt.Sprintf(template, args...))
}

func Printf(format string, args ...interface{}) {
	Debug(format, args...)
}

// Flush function call to clear logger buffer
func Flush() {
	inst.Info("Flushing log buffer")
	_ = inst.Sync()
}
