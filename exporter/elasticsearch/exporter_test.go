package elasticsearch

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/k8shuginn/event-collector/dummy"
	"github.com/k8shuginn/event-collector/pkg/logger"
)

const (
	envEsAddrs = "ELASTICSEARCH_ADDR"
	envEsIndex = "ELASTICSEARCH_INDEX"
)

func TestExporter(t *testing.T) {
	logger.CreateGlobalLogger("es_test")

	// get config from env
	esAddr := os.Getenv(envEsAddrs)
	if esAddr == "" {
		t.Fatalf("elasticsearch env %s not set, use default", envEsAddrs)
	}

	esIndex := os.Getenv(envEsIndex)
	if esIndex == "" {
		t.Fatalf("elasticsearch env %s not set, use default", envEsAddrs)
	}

	// run exporter
	es, err := NewElasticsearchExporter(
		[]string{esAddr}, esIndex,
	)
	if err != nil {
		t.Fatal(err)
	}

	// write test data
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
