package volume

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/k8shuginn/event-collector/dummy"
	"github.com/k8shuginn/event-collector/pkg/logger"
	"github.com/k8shuginn/event-collector/pkg/util"
)

func TestExporter(t *testing.T) {
	logger.CreateGlobalLogger("volume_test")

	// get current path
	path, _ := util.GetRootPath()
	dir := filepath.Join(path, "test")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// create volume exporter
	exporter, err := NewVolumeExporter(
		"testevent", dir,
		WithMaxFileCount(5),
		WithMaxFileSize(1024),
	)
	if err != nil {
		t.Fatal(err)
	}

	// write dummy data
	dummy := dummy.MakeDummy("1")
	dummyBytes, _ := json.Marshal(dummy)
	for i := 0; i < 10; i++ {
		if err := exporter.writeData(dummyBytes); err != nil {
			t.Fatal(err)
		}
	}
}
