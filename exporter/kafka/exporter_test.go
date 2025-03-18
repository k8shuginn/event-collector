package kafka

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/k8shuginn/event-collector/dummy"
	"github.com/k8shuginn/event-collector/pkg/logger"
)

func TestExporte(t *testing.T) {
	logger.CreateGlobalTestLogger()

	kafka, err := NewKafkaExporter([]string{"localhost:31719"}, "event")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go kafka.Start(ctx)

	eventData := dummy.MakeDummy("1")
	eventBytes, _ := json.Marshal(eventData)
	kafka.Write(eventBytes)

	eventData = dummy.MakeDummy("2")
	eventBytes, _ = json.Marshal(eventData)
	kafka.Write(eventBytes)

	cancel()
}
