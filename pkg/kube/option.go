package kube

import (
	"strings"
	"time"
)

type config struct {
	kubeConfig string
	resyncTime time.Duration

	namespaces map[string]struct{}
}

func defaultConfig() *config {
	return &config{
		kubeConfig: "",
		resyncTime: 0,
		namespaces: make(map[string]struct{}),
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

// WithNamespaces 수집할 namespace 리스트 설정
// namespace의 값이 없는 경우 전체 네임스페이스를 모니터링합니다.
func WithNamespaces(namespaces []string) Option {
	return func(c *config) {
		if len(namespaces) > 0 {
			for _, ns := range namespaces {
				if ns == "" {
					continue
				}

				c.namespaces[strings.ToLower(ns)] = struct{}{}
			}
		}
	}
}
