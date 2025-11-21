package kafka

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/k8shuginn/event-collector/dummy"
	"github.com/k8shuginn/event-collector/pkg/logger"
)

const (
	envKafkaBrokers = "KAFKA_BROKERS"
	envKafkaTopic   = "KAFKA_TOPIC"
)

func TestExporte(t *testing.T) {
	logger.CreateGlobalLogger("kafka_test")

	// get config from env
	kafkaBrokers := os.Getenv(envKafkaBrokers)
	if kafkaBrokers == "" {
		t.Fatalf("env %s not set, use default", envKafkaBrokers)
	}

	kafkaTopic := os.Getenv(envKafkaTopic)
	if kafkaTopic == "" {
		t.Fatalf("env %s not set, use default", envKafkaTopic)
	}

	// run exporter
	kafka, err := NewKafkaExporter([]string{kafkaBrokers}, kafkaTopic)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	go kafka.Start(ctx, &wg)

	// write test data
	eventData := dummy.MakeDummy("1")
	eventBytes, _ := json.Marshal(eventData)
	kafka.Write(eventBytes)

	eventData = dummy.MakeDummy("2")
	eventBytes, _ = json.Marshal(eventData)
	kafka.Write(eventBytes)

	// stop exporter
	cancel()
}
