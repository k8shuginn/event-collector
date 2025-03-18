package elasticsearch

import "time"

type config struct {
	user string
	pass string

	chanSize  int
	flushTime time.Duration
	flushSize int
}

type Option func(*config)

func fromOptions(opts ...Option) *config {
	c := &config{
		chanSize:  200,
		flushTime: 1 * time.Second, // 1초
		flushSize: 1024 * 1024,     // 1MB
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithUser elasticsearch user 설정
func WithUser(user string) Option {
	return func(c *config) {
		if user != "" {
			c.user = user
		}
	}
}

// WithPass elasticsearch password 설정
func WithPass(pass string) Option {
	return func(c *config) {
		if pass != "" {
			c.pass = pass
		}
	}
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
func WithFlushTime(d time.Duration) Option {
	return func(c *config) {
		if d > 0 {
			c.flushTime = d
		}
	}
}

// WithFlushSize elasticsearch exporter buffer flush 사이즈 설정
func WithFlushSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.flushSize = size
		}
	}
}
