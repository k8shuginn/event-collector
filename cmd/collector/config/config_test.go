package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/k8shuginn/event-collector/pkg/logger"
)

func TestConfig(t *testing.T) {
	logger.CreateGlobalTestLogger()
	path, _ := os.Getwd()

	config, err := LoadConfig(filepath.Join(path, "test.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(config)
}
