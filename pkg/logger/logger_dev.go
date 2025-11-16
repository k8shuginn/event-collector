//go:build !prd
// +build !prd

package logger

import (
	"go.uber.org/zap"
)

// CreateGlobalLogger 전역 로거 생성
// appName: 앱 이름
// opts: logger option
func CreateGlobalLogger(appName string, options ...Option) error {
	_ = fromOptions(appName, options...)
	writer = zap.NewNop()

	return nil
}

// CreatLogger 로거 생성
// appName: 앱 이름
// opts: logger option
func CreatLogger(appName string, options ...Option) (*zap.Logger, error) {
	_ = fromOptions(appName, options...)
	nopLogger := zap.NewNop()

	return nopLogger, nil
}
