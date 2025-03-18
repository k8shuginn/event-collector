package volume

type config struct {
	maxFileSize  int
	maxFileCount int
	chanSize     int
}

type Option func(*config)

func defaultConfig() *config {
	return &config{
		maxFileSize:  50 * 1024 * 1024, // 50MB
		maxFileCount: 10,               // 10 files
		chanSize:     200,
	}
}

func fromOptions(options ...Option) *config {
	c := defaultConfig()
	for _, option := range options {
		option(c)
	}

	return c
}

// WithMaxFileSize 최대 파일 크기 설정
func WithMaxFileSize(size int) Option {
	return func(c *config) {
		if size > 0 {
			c.maxFileSize = size
		}
	}
}

// WithMaxFileCount 최대 파일 개수 설정
func WithMaxFileCount(count int) Option {
	return func(c *config) {
		if count > 0 {
			c.maxFileCount = count
		}
	}
}
