package elasticsearch

import (
	"encoding/json"
	"testing"

	"github.com/k8shuginn/event-collector/dummy"
	"github.com/k8shuginn/event-collector/pkg/logger"
)

func TestExporter(t *testing.T) {
	logger.CreateGlobalTestLogger()

	es, err := NewElasticsearchExporter(
		[]string{"https://localhost:30090"}, "event",
		WithUser("elastic"),
		WithPass("elastic"),
	)
	if err != nil {
		t.Fatal(err)
	}

	eventData := dummy.MakeDummy("1")
	eventBytes, _ := json.Marshal(eventData)
	es.writeBuffer(eventBytes)

	eventData = dummy.MakeDummy("2")
	eventBytes, _ = json.Marshal(eventData)
	es.writeBuffer(eventBytes)

	eventData = dummy.MakeDummy("3")
	eventBytes, _ = json.Marshal(eventData)
	es.writeBuffer(eventBytes)

	es.writeBulk()
}
