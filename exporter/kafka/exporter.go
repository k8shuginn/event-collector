package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/k8shuginn/event-collector/exporter"
	"github.com/k8shuginn/event-collector/pkg/logger"
	"go.uber.org/zap"
)

var _ exporter.Exporter = (*KafkaExporter)(nil)

type KafkaExporter struct {
	producer sarama.AsyncProducer
	topic    string
}

// NewKafkaExporter kafka exporter 생성
// brokers: kafka brokers
// topic: kafka topic
// opts: exporter option
func NewKafkaExporter(brokers []string, topic string, opts ...Option) (*KafkaExporter, error) {
	c := fromOptions(opts...)

	producer, err := sarama.NewAsyncProducer(brokers, c.saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka exporter: %w", err)
	}

	return &KafkaExporter{
		producer: producer,
		topic:    topic,
	}, nil
}

// Start exporter 시작
func (e *KafkaExporter) Start(ctx context.Context) error {
	logger.Info("[kafka exporter] started")
	defer func() {
		e.shutdown()
		logger.Info("[kafka exporter] stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-e.producer.Errors():
			logger.Error("[kafka exporter] failed to send message", zap.Error(err), zap.String("topic", err.Msg.Topic), zap.Int32("partition", err.Msg.Partition))
		case success := <-e.producer.Successes():
			logger.Debug("[kafka exporter] message sent", zap.String("topic", success.Topic), zap.Int32("partition", success.Partition), zap.Int64("offset", success.Offset))
		}
	}
}

// shutdown exporter 종료
func (e *KafkaExporter) shutdown() {
	e.producer.AsyncClose()
}

// Write 데이터 전송
func (e *KafkaExporter) Write(data []byte) {
	e.producer.Input() <- &sarama.ProducerMessage{
		Topic: e.topic,
		Key:   nil,
		Value: sarama.ByteEncoder(data),
	}
}
