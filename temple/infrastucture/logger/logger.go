package logger

import (
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option = func(config *Config)

type Config struct {
	Encoding string
	Level    Level
	DevMode  bool
	AppName  string
}

func (c *Config) GetLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[c.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

type Level string

const (
	Debug  Level = "debug"
	Info   Level = "info"
	Warn   Level = "warn"
	Error  Level = "error"
	DPanic Level = "dpanic"
	Panic  Level = "panic"
	Fatal  Level = "fatal"
)

var loggerLevelMap = map[Level]zapcore.Level{
	Debug:  zapcore.DebugLevel,
	Info:   zapcore.InfoLevel,
	Warn:   zapcore.WarnLevel,
	Error:  zapcore.ErrorLevel,
	DPanic: zapcore.DPanicLevel,
	Panic:  zapcore.PanicLevel,
	Fatal:  zapcore.FatalLevel,
}

func (c *Config) Use(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	logLevel := c.GetLoggerLevel()
	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	if c.DevMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.NameKey = "service"
	encoderCfg.TimeKey = "time"
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "line"
	encoderCfg.MessageKey = "message"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	if c.Encoding == "console" {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		encoderCfg.ConsoleSeparator = " | "
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoderCfg.FunctionKey = "caller"
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))

	lg := zap.New(core)
	if c.DevMode {
		lg = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	codeOtelZap := otelzap.New(
		lg.Named(os.Getenv("TECH_SERVICE_NAME")),
		otelzap.WithTraceIDField(true),
		otelzap.WithMinLevel(logLevel),
		otelzap.WithStackTrace(c.DevMode),
	)
	otelzap.ReplaceGlobals(codeOtelZap)
}
