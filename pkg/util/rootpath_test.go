package util

import (
	"os"
	"testing"
)

// TestGetRootPath 루트 경로를 올바르게 찾는지 테스트합니다.
func TestGetRootPath(t *testing.T) {
	path, err := GetRootPath()
	if err != nil {
		t.Errorf("failed to get root path: %v", err)
	}

	fInfo, err := os.Stat(path)
	if err != nil {
		t.Errorf("failed to stat root path: %v", err)
	}

	if !fInfo.IsDir() {
		t.Errorf("expected root path to be a directory, got file")
	}

	t.Logf("project root path: %s", path)
}
