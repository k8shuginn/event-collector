package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	path, _ := os.Getwd()

	config, err := LoadConfig(filepath.Join(path, "test.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(config)
}
