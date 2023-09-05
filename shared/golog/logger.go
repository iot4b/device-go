package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

type Logger struct {
	level Level
	zap   *zap.SugaredLogger
}
type Level zapcore.Level

type Zap struct {
	sugarClient *zap.SugaredLogger
	client      *zap.Logger
}

var l = &Logger{}

func Init(debug bool) {
	lvl := "info"
	isDev := false
	disableStack := true

	// setup logs
	if debug {
		lvl = "debug"
		isDev = true
		disableStack = false
	}

	config := &zap.Config{
		Level:             LevelToAtomic(ParseLevel(lvl)),
		Development:       isDev,
		DisableCaller:     false,
		DisableStacktrace: disableStack,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     stampTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   callerEncoder,
		},
		//OutputPaths:      []string{"/var/log/syslog"},
		//ErrorOutputPaths: []string{"/var/log/syslog"},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	lg, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		fmt.Println("Logger init error: ", err)
		return
	}

	l = &Logger{
		level: ParseLevel(lvl),
		zap:   lg.Sugar(),
	}
}

func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	arr := strings.Split(caller.Function, ".")
	funName := arr[len(arr)-1]
	enc.AppendString(caller.TrimmedPath() + "." + funName + "()")
}

func stampTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	format := "Jan 02 15:04:05"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, format)
		return
	}

	enc.AppendString(t.Format(format))
}

func LevelToAtomic(lvl Level) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(zapcore.Level(lvl))
}

func ParseLevel(lvl string) (level Level) {
	switch lvl {
	case "debug":
		level = Level(zap.DebugLevel)
	case "info":
		level = Level(zap.InfoLevel)
	case "warning":
		level = Level(zap.WarnLevel)
	case "error":
		level = Level(zap.ErrorLevel)
	case "panic":
		level = Level(zap.PanicLevel)
	case "fatal":
		level = Level(zap.FatalLevel)
	}
	return
}

// Fatal followed by a call to os.Exit(1).
func Fatal(msg ...interface{}) {
	l.zap.Fatal(msg)
	os.Exit(1)
}

// Fatalf followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.zap.Fatal(msg)
	os.Exit(1)
}

// Panic followed by a call to panic().
func Panic(msg ...interface{}) {
	l.zap.Panic(msg)
}

// Panicf followed by a call to panic().
func Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.zap.Panic(msg)
}

// Error logs a message using ERROR as log level.
func Error(msg ...interface{}) {
	l.zap.Error(msg)
}

// Errorf logs a message using ERROR as log level.
func Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.zap.Error(msg)
}

// Warning logs a message using WARNING as log level.
func Warning(msg ...interface{}) {
	l.zap.Warn(msg)
}

// Warningf logs a message using WARNING as log level.
func Warningf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.zap.Warn(msg)
}

// Info logs a message using INFO as log level.
func Info(msg ...interface{}) {
	l.zap.Info(msg)
}

// Infof logs a message using INFO as log level.
func Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.zap.Info(msg)
}

// Debug logs a message using DEBUG as log level.
func Debug(msg ...interface{}) {
	l.zap.Debug(msg)
}

// Debugf logs a message using DEBUG as log level.
func Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.zap.Debug(msg)
}

// Fatalw followed by a call to os.Exit(1).
func Fatalw(msg string, args ...interface{}) {
	l.zap.Fatalw(msg, args...)
	os.Exit(1)
}

// Errorw logs a message using ERROR as log level.
func Errorw(msg string, args ...interface{}) {
	l.zap.Errorw(msg, args...)
}

// Warningf logs a message using WARNING as log level.
func Warningw(msg string, args ...interface{}) {
	l.zap.Warnw(msg, args...)
}

// Infof logs a message using INFO as log level.
func Infow(msg string, args ...interface{}) {
	l.zap.Infow(msg, args...)
}

// Debugf logs a message using DEBUG as log level.
func Debugw(msg string, args ...interface{}) {
	l.zap.Debugw(msg, args...)
}
