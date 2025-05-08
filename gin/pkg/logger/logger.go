package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// log encoding
const (
	JSONEncode    string = "json"
	ConsoleEncode string = "console"
)

// log level
const (
	DebugLevel string = "debug"
	InfoLevel  string = "info"
	WarnLevel  string = "warn"
	ErrorLevel string = "error"
	FatalLevel string = "fatal"
	PanicLevel string = "panic"
)

var levelMapping = map[string]zap.AtomicLevel{
	"debug": zap.NewAtomicLevelAt(zap.DebugLevel),
	"info":  zap.NewAtomicLevelAt(zap.InfoLevel),
	"warn":  zap.NewAtomicLevelAt(zap.WarnLevel),
	"error": zap.NewAtomicLevelAt(zap.ErrorLevel),
	"fatal": zap.NewAtomicLevelAt(zap.FatalLevel),
	"panic": zap.NewAtomicLevelAt(zap.PanicLevel),
}

type Config struct {
	Encoding string `json:"encoding" yaml:"encoding"`
	Level    string `json:"level" yaml:"level"`
}

var logger *zap.Logger

func customColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch level {
	case zapcore.DebugLevel:
		color = "\033[36m" // 青色
	case zapcore.InfoLevel:
		color = "\033[32m" // 绿色
	case zapcore.WarnLevel:
		color = "\033[33m" // 黄色
	case zapcore.ErrorLevel:
		color = "\033[31m" // 红色
	case zapcore.FatalLevel, zapcore.PanicLevel:
		color = "\033[35m" // 紫色
	default:
		color = "\033[0m" // 默认
	}
	enc.AppendString(color + level.String() + "\033[0m")
}

func Initialize(encoding string, level string) {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		CallerKey:     "caller",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if encoding == JSONEncode {
		encoder = zapcore.NewJSONEncoder(config)
	}

	if encoding == ConsoleEncode {
		config.EncodeLevel = customColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(config)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), levelMapping[level])

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Info(msg string) {
	logger.Sugar().Info(msg)
}

func Infof(format string, args ...any) {
	logger.Sugar().Infof(format, args...)
}

func Debug(msg string) {
	logger.Sugar().Debug(msg)
}

func Debugf(format string, args ...any) {
	logger.Sugar().Debugf(format, args...)
}

func Warn(msg string) {
	logger.Sugar().Warn(msg)
}

func Warnf(format string, args ...any) {
	logger.Sugar().Warnf(format, args...)
}

func Error(msg string) {
	logger.Sugar().Error(msg)
}

func Errorf(format string, args ...any) {
	logger.Sugar().Errorf(format, args...)
}

func Fatal(msg string) {
	logger.Sugar().Fatal(msg)
}

func Fatalf(format string, args ...any) {
	logger.Sugar().Fatalf(format, args...)
}

func Panic(msg string) {
	logger.Sugar().Panic(msg)
}
