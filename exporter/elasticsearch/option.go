package elasticsearch

import (
	"os"
	"time"
)

const (
	EnvElaticsearchUser = "ELASTICSEARCH_USER"
	EnvElaticsearchPass = "ELASTICSEARCH_PASSWORD"
)

type config struct {
	user, pass string
	chanSize   int
	flushTime  time.Duration
	flushSize  int
}

type Option func(*config)

func defaultConfig() *config {
	return &config{
		user:      os.Getenv(EnvElaticsearchUser),
		pass:      os.Getenv(EnvElaticsearchPass),
		chanSize:  200,
		flushTime: 1 * time.Second, // 1초
		flushSize: 1024 * 1024,     // 1MB
	}
}

func fromOptions(opts ...Option) *config {
	c := defaultConfig()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithChanSize elasticsearch exporter buffer channel size 설정
func WithChanSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.chanSize = size
		}
	}
}

// WithFlushTime elasticsearch exporter buffer flush 시간 설정
func WithFlushTime(seconds int) Option {
	return func(c *config) {
		if seconds > 0 {
			c.flushTime = time.Duration(seconds) * time.Second
		}
	}
}

// WithFlushSize elasticsearch exporter buffer flush 사이즈 설정
func WithFlushSize(byteSize int) Option {
	return func(c *config) {
		if byteSize > 0 {
			c.flushSize = byteSize
		}
	}
}
