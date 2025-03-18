package logger

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func Test_logger(t *testing.T) {
	nowPath, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	if err := CreateGlobalLogger(
		"testlog",
		WithPath(nowPath),
		WithLogLevel("debug"),
		WithLogLocalTime(false),
		WithLogCompress(true),
		WitchEncoder("console"),
	); err != nil {
		t.Error(err)
	}

	writer.Info("info message")
	writer.Debug("debug message", zap.String("key", "value"))
	writer.Warn("warn message", zap.String("key", "value"))
	writer.Error("error message", zap.String("key", "value"))

	// 현재 디렉토리에서 패턴과 일치하는 파일 검색
	pattern := "testlog*.log"
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Error(err)
	}

	// 파일이 없는 경우 처리
	if len(files) == 0 {
		t.Log("not found log file")
		return
	}

	// 파일 삭제
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			t.Log("failed to remove file: ", file)
		} else {
			t.Log("remove file: ", file)
		}
	}
}
