//go:build prd
// +build prd

package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateGlobalLogger 전역 로거 생성
// appName: 앱 이름
// opts: logger option
func CreateGlobalLogger(appName string, options ...Option) error {
	c := fromOptions(appName, options...)
	writer = zap.New(
		zapcore.NewCore(
			c.encoder,
			zapcore.NewMultiWriteSyncer(append([]zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}, zapcore.AddSync(&c.logger))...),
			c.level,
		),
	)

	if writer == nil {
		return fmt.Errorf("failed to create logger")
	}

	return nil
}

// CreatLogger 로거 생성
// appName: 앱 이름
// opts: logger option
func CreatLogger(appName string, options ...Option) (*zap.Logger, error) {
	c := fromOptions(appName, options...)
	logger := zap.New(
		zapcore.NewCore(
			c.encoder,
			zapcore.NewMultiWriteSyncer(append([]zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}, zapcore.AddSync(&c.logger))...),
			c.level,
		),
	)

	if logger == nil {
		return nil, fmt.Errorf("failed to create logger")
	}

	return logger, nil
}
