package volume

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/k8shuginn/event-collector/dummy"
	"github.com/k8shuginn/event-collector/pkg/logger"
)

func TestExporter(t *testing.T) {
	logger.CreateGlobalTestLogger()

	// get current path
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	exporter, err := NewVolumeExporter(
		"testevent", path,
		WithMaxFileCount(5),
		WithMaxFileSize(1024),
	)
	if err != nil {
		t.Fatal(err)
	}

	dummy := dummy.MakeDummy("1")
	dummyBytes, _ := json.Marshal(dummy)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
	exporter.writeData(dummyBytes)
}
