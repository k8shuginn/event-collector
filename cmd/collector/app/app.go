package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/k8shuginn/event-collector/cmd/collector/config"
	"github.com/k8shuginn/event-collector/exporter"
	"github.com/k8shuginn/event-collector/exporter/elasticsearch"
	"github.com/k8shuginn/event-collector/exporter/kafka"
	"github.com/k8shuginn/event-collector/exporter/volume"
	"github.com/k8shuginn/event-collector/pkg/kube"
	"github.com/k8shuginn/event-collector/pkg/logger"
)

const (
	DefaultTimeout = 10 * time.Second
)

type Collector struct {
	k8sClient *kube.Client
	handler   *Handler

	exporters []exporter.Exporter
}

// NewCollector 설정에 따라 Collector 생성
func NewCollector(c *config.Config) (*Collector, error) {
	exporters, err := createExporters(c)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		exporters: exporters,
	}

	client, err := kube.NewClient(handler,
		kube.WithKubeConfig(c.Kube.Config),
		kube.WithResycPeriod(c.Kube.Resync),
		kube.WithNamespaces(c.Kube.Namespaces),
	)
	if err != nil {
		return nil, err
	}

	return &Collector{
		k8sClient: client,
		handler:   handler,
		exporters: exporters,
	}, nil
}

// Run Collector 실행
func (c *Collector) Run() {
	logger.Info("kubernetes event collector started ...")
	defer logger.Info("kubernetes event collector stopped ...")

	wg := sync.WaitGroup{}
	// 시그널 수신 설정
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// event 수집 시작
	ctx, cancel := context.WithCancel(context.Background())
	for _, e := range c.exporters {
		go e.Start(ctx, &wg)
	}
	c.k8sClient.Run()

	// 시그널 수신 시 종료
	<-sigChan
	c.k8sClient.Close()
	cancel()

	// exporter 종료 대기 (최대 10초)
	exitChan := make(chan struct{})
	go func() {
		defer close(exitChan)
		wg.Wait()
	}()
	select {
	case <-exitChan:
		logger.Info("all exporters stopped")
		return
	case <-time.After(DefaultTimeout):
		logger.Warn("timeout waiting for exporters to stop")
		return
	}
}

// createExporters 설정에 따라 Exporter 생성
func createExporters(c *config.Config) ([]exporter.Exporter, error) {
	exporters := []exporter.Exporter{}

	// Kafka Exporter 생성
	if c.Kafka.Enable {
		kafkaExporter, err := kafka.NewKafkaExporter(
			c.Kafka.Brokers, c.Kafka.Topic,
			kafka.WithTimeout(c.Kafka.Timeout),
			kafka.WithRetry(c.Kafka.Retry),
			kafka.WithRetryBackoff(c.Kafka.RetryBackoff),
			kafka.WithFlushMsg(c.Kafka.FlushMsg),
			kafka.WithFlushTime(c.Kafka.FlushTime),
			kafka.WithFlushByte(c.Kafka.FlushByte),
		)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, kafkaExporter)
	}

	// ElasticSearch Exporter 생성
	if c.ElasticSearch.Enable {
		elasticExporter, err := elasticsearch.NewElasticsearchExporter(
			c.ElasticSearch.Addresses, c.ElasticSearch.Index,
			elasticsearch.WithChanSize(c.ElasticSearch.ChanSize),
			elasticsearch.WithFlushTime(c.ElasticSearch.FlushTime),
			elasticsearch.WithFlushSize(c.ElasticSearch.FlushSize),
		)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, elasticExporter)
	}

	// Volume Exporter 생성
	if c.Volume.Enable {
		volumeExporter, err := volume.NewVolumeExporter(
			c.Volume.FileName, c.Volume.FilePath,
			volume.WithMaxFileSize(c.Volume.MaxFileSize),
			volume.WithMaxFileCount(c.Volume.MaxFileCount),
		)
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, volumeExporter)
	}

	return exporters, nil
}
