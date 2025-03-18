package kube

import "time"

type config struct {
	kubeConfig string
	resyncTime time.Duration
}

func defaultConfig() *config {
	return &config{
		kubeConfig: "",
		resyncTime: 0,
	}
}

type Option func(*config)

func fromOptions(options ...Option) *config {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}
	return config
}

// WithResyncTime resync 시간 설정
func WithResyncTime(resyncTime time.Duration) Option {
	return func(c *config) {
		if resyncTime > 0 {
			c.resyncTime = resyncTime
		}
	}
}

// WithKubeConfig kubeconfig 설정
func WithKubeConfig(kubeConfig string) Option {
	return func(c *config) {
		if kubeConfig != "" {
			c.kubeConfig = kubeConfig
		}
	}
}
