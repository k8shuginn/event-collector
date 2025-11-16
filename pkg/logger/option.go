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

	DefaultMaxSize    int  = 100   // megabytes
	DefaultMaxBackups int  = 3     // number of backups
	DefaultMaxAge     int  = 7     // log retention period days
	DefaultLocalTime  bool = true  // use local time for timestamps
	DefaultCompress   bool = false // disabled compression by default
)

type EncoderType string

const (
	JSONEncoder    EncoderType = "json"
	ConsoleEncoder EncoderType = "console"
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
func WithPath(dirPath string) Option {
	return func(c *config) {
		if dirPath != "" {
			c.logger.Filename = filepath.Join(dirPath, c.appName+DefaultLogExtention)
		}
	}
}

// WithLogExtention 로그 파일 확장자 설정
func WithLogMaxSize(maxSize int) Option {
	return func(c *config) {
		if maxSize > 0 {
			c.logger.MaxSize = maxSize
		}
	}
}

// WithLogMaxBackups 로그 파일 백업 설정
func WithLogMaxBackups(maxBackups int) Option {
	return func(c *config) {
		if maxBackups > 0 {
			c.logger.MaxBackups = maxBackups
		}
	}
}

// WithLogMaxAge 로그 파일 보관 기간 설정
func WithLogMaxAge(maxAge int) Option {
	return func(c *config) {
		if maxAge > 0 {
			c.logger.MaxAge = maxAge
		}
	}
}

// WithLogLocalTime 로컬 시간 설정
func WithLogLocalTime(isLocalTime bool) Option {
	return func(c *config) {
		c.logger.LocalTime = isLocalTime
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

// WithEncoder 로그 인코더 설정
// JSON, Console 지원
func WithEncoder(encoderType EncoderType) Option {
	return func(c *config) {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		switch encoderType {
		case ConsoleEncoder:
			c.encoder = zapcore.NewConsoleEncoder(encoderConfig)
		case JSONEncoder:
			fallthrough
		default:
			c.encoder = zapcore.NewJSONEncoder(encoderConfig)
		}
	}
}
