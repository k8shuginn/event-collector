package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/k8shuginn/event-collector/exporter"
	"github.com/k8shuginn/event-collector/pkg/logger"
	"go.uber.org/zap"
)

var _ exporter.Exporter = (*ElasticsearchExporter)(nil)

type ElasticsearchExporter struct {
	es    *elasticsearch.Client
	index string

	// buffer
	dataChan  chan []byte
	mux       *sync.Mutex
	flushTime time.Duration
	flushSize int
	buffer    []byte
}

// NewElasticsearchExporter elasticsearch exporter 생성
// addrs: elasticsearch 주소
// index: elasticsearch index
// opts: exporter option
func NewElasticsearchExporter(addrs []string, index string, opts ...Option) (*ElasticsearchExporter, error) {
	c := fromOptions(opts...)

	esCfg := elasticsearch.Config{
		Addresses: addrs,
		Username:  c.user,
		Password:  c.pass,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // TLS 인증서 검증 비활성화
			},
		},
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch exporter: %w", err)
	}

	_, err = es.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}

	e := &ElasticsearchExporter{
		es:        es,
		index:     index,
		dataChan:  make(chan []byte, c.chanSize),
		mux:       &sync.Mutex{},
		flushTime: c.flushTime,
		flushSize: c.flushSize,
		buffer:    make([]byte, 0),
	}

	return e, nil
}

// Start exporter 시작
func (e *ElasticsearchExporter) Start(ctx context.Context, wg *sync.WaitGroup) error {
	logger.Info("[elasticsearch exporter] started")
	ticker := time.NewTicker(5 * time.Second)

	// shutdown 처리
	defer func() {
		close(e.dataChan)
		ticker.Stop()
		e.shutdown()

		logger.Info("[elasticsearch exporter] stopped")
		wg.Done()
	}()

	wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-e.dataChan:
			e.writeBuffer(data)
			if len(e.buffer) >= e.flushSize {
				e.writeBulk()
			}
		case <-ticker.C:
			if len(e.buffer) > 0 {
				e.writeBulk()
			}
		}
	}
}

// writeBulk bulk flush
func (e *ElasticsearchExporter) writeBulk() {
	e.mux.Lock()
	defer e.mux.Unlock()

	buf := bytes.NewBuffer(e.buffer)
	res, err := e.es.Bulk(buf, e.es.Bulk.WithContext(context.Background()))
	if err != nil {
		logger.Error("[elasticsearch exporter] failed to write bulk", zap.Error(err))
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Error("[elasticsearch exporter] bulk request failed", zap.String("response", res.String()))
		return
	}

	logger.Debug("[elasticsearch exporter] bulk request success")

	// buffer 초기화
	e.buffer = e.buffer[:0]
}

// writeBuffer buffer write
func (e *ElasticsearchExporter) writeBuffer(data []byte) {
	e.mux.Lock()
	defer e.mux.Unlock()

	result := make([]byte, 0)
	meta := fmt.Sprintf(`{"index": {"_index": "%s" } }`, e.index)
	result = append(result, meta...)
	result = append(result, '\n')
	result = append(result, data...)
	result = append(result, '\n')

	e.buffer = append(e.buffer, result...)
}

// Write data write
func (e *ElasticsearchExporter) Write(data []byte) {
	e.dataChan <- data
}

func (e *ElasticsearchExporter) shutdown() {
	e.writeBulk()
}
