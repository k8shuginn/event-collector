package logger

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultPath         string = "/var/log/app/"
	DefaultLogExtention string = ".log"
	DefaultAppName      string = "app"

	DefaultMaxSize    int  = 100
	DefaultMaxBackups int  = 3
	DefaultMaxAge     int  = 7
	DefaultLocalTime  bool = true
	DefaultCompress   bool = true
)

type config struct {
	appName string
	encoder zapcore.Encoder
	level   zapcore.Level
	logger  lumberjack.Logger
}

// defaultOption 기본 설정
func defaultOption() config {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return config{
		appName: DefaultAppName,
		encoder: zapcore.NewJSONEncoder(encoderConfig),
		level:   zapcore.InfoLevel,
		logger: lumberjack.Logger{
			Filename:   DefaultPath + DefaultAppName + DefaultLogExtention,
			MaxSize:    DefaultMaxSize,
			MaxBackups: DefaultMaxBackups,
			MaxAge:     DefaultMaxAge,
			LocalTime:  DefaultLocalTime,
			Compress:   DefaultCompress,
		},
	}
}

type Option func(*config)

// fromOptions 설정 적용
func fromOptions(appName string, options ...Option) *config {
	c := defaultOption()
	c.appName = appName

	for _, option := range options {
		option(&c)
	}

	return &c
}

// WithPath 로그 파일 경로 설정
func WithPath(path string) Option {
	return func(c *config) {
		if path != "" {
			c.logger.Filename = filepath.Join(path, c.appName+DefaultLogExtention)
		}
	}
}

// WithLogExtention 로그 파일 확장자 설정
func WithLogMaxSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.logger.MaxSize = size
		}
	}
}

// WithLogMaxBackups 로그 파일 백업 설정
func WithLogMaxBackups(backups int) Option {
	return func(c *config) {
		if backups > 0 {
			c.logger.MaxBackups = backups
		}
	}
}

// WithLogMaxAge 로그 파일 보관 기간 설정
func WithLogMaxAge(age int) Option {
	return func(c *config) {
		if age > 0 {
			c.logger.MaxAge = age
		}
	}
}

// WithLogLocalTime 로컬 시간 설정
func WithLogLocalTime(localTime bool) Option {
	return func(c *config) {
		c.logger.LocalTime = localTime
	}
}

// WithLogCompress 로그 압축 설정
func WithLogCompress(compress bool) Option {
	return func(c *config) {
		c.logger.Compress = compress
	}
}

// WithLogLevel 로그 레벨 설정
func WithLogLevel(level string) Option {
	return func(c *config) {
		switch strings.ToUpper(level) {
		case "DEBUG":
			c.level = zapcore.DebugLevel
		case "WARN", "WARNING":
			c.level = zapcore.WarnLevel
		case "ERROR", "ERR":
			c.level = zapcore.ErrorLevel
		case "DPANIC":
			c.level = zapcore.DPanicLevel
		case "PANIC":
			c.level = zapcore.PanicLevel
		case "FATAL":
			c.level = zapcore.FatalLevel
		case "INFO", "INF":
			fallthrough
		default:
			c.level = zapcore.InfoLevel
		}
	}
}

// WitchEncoder 로그 인코더 설정
func WitchEncoder(encoder string) Option {
	return func(c *config) {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		switch strings.ToUpper(encoder) {
		case "CONSOLE":
			c.encoder = zapcore.NewConsoleEncoder(encoderConfig)
		case "JSON":
			fallthrough
		default:
			c.encoder = zapcore.NewJSONEncoder(encoderConfig)
		}
	}
}
