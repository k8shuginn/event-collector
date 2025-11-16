package kube

import (
	"strings"
	"time"
)

const (
	defaultKubeConfig   = ""
	defaultResyncPeriod = 0 * time.Minute // no resync by default
	minResyncPeriod     = 10 * time.Minute
)

type config struct {
	kubeConfig   string
	resyncPeriod time.Duration
	namespaces   map[string]struct{}
}

func defaultConfig() *config {
	return &config{
		kubeConfig:   defaultKubeConfig,
		resyncPeriod: defaultResyncPeriod,
		namespaces:   make(map[string]struct{}),
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

// WithKubeConfig kubeconfig 설정
func WithKubeConfig(kubeConfig string) Option {
	return func(c *config) {
		if kubeConfig != "" {
			c.kubeConfig = kubeConfig
		}
	}
}

// WithResycPeriod resync 시간 설정
func WithResycPeriod(resyncPeriod time.Duration) Option {
	return func(c *config) {
		if resyncPeriod <= 0 {
			c.resyncPeriod = 0
		} else if 0 < resyncPeriod && resyncPeriod < minResyncPeriod {
			c.resyncPeriod = minResyncPeriod
		} else {
			c.resyncPeriod = resyncPeriod
		}
	}
}

// WithNamespaces 수집할 namespace 리스트 설정
// namespace의 값이 없는 경우 전체 네임스페이스를 모니터링합니다.
func WithNamespaces(namespaces []string) Option {
	return func(c *config) {
		if len(namespaces) == 0 {
			return
		}

		c.namespaces = make(map[string]struct{})
		for _, ns := range namespaces {
			if ns == "" {
				continue
			}

			c.namespaces[strings.ToLower(ns)] = struct{}{}
		}
	}
}
