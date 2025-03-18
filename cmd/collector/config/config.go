package config

import (
	"fmt"
	"os"
	"time"

	"github.com/k8shuginn/event-collector/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Kube struct {
		Config string        `yaml:"config"` // 없으면 in-cluster 자동 설정
		Resync time.Duration `yaml:"resync"`
	} `yaml:"kube"`

	Kafka struct {
		Enable       bool          `yaml:"enable"`  // 필수
		Brokers      []string      `yaml:"brokers"` // 필수
		Topic        string        `yaml:"topic"`   // 필수
		Timeout      time.Duration `yaml:"timeout"`
		Retry        int           `yaml:"retry"`
		RetryBackoff time.Duration `yaml:"retryBackoff"`
		FlushMsg     int           `yaml:"flushMsg"`
		FlushTime    time.Duration `yaml:"flushTime"`
		FlushByte    int           `yaml:"flushByte"`
	} `yaml:"kafka"`

	ElasticSearch struct {
		Enable    bool     `yaml:"enable"`    // 필수
		Addresses []string `yaml:"addresses"` // 필수
		User      string   `yaml:"user"`      // 필수
		Pass      string   `yaml:"pass"`      // 필수
		Index     string   `yaml:"index"`     // 필수
	} `yaml:"elasticsearch"`

	Volume struct {
		Enable       bool   `yaml:"enable"`   // 필수
		FileName     string `yaml:"fileName"` // 필수
		FilePath     string `yaml:"filePath"` // 필수
		MaxFileSize  int    `yaml:"maxFileSize"`
		MaxFileCount int    `yaml:"maxFileCount"`
	} `yaml:"volume"`
}

// LoadConfig 설정 파일을 읽어서 Config 구조체로 반환
func LoadConfig(fileName string) (*Config, error) {
	config := &Config{}
	if err := readFile(fileName, config); err != nil {
		return nil, fmt.Errorf("failed config read file: %w", err)
	}
	showConfig(config)

	return config, checkConfig(config)
}

// readFile 설정 파일을 읽어서 Config 구조체로 반환
func readFile(fileName string, config *Config) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return err
	}

	return nil
}

// showConfig 설정 파일을 로그로 출력
func showConfig(config *Config) {
	logger.Debug("kubernetes config",
		zap.String("config", config.Kube.Config),
		zap.Duration("resync", config.Kube.Resync),
	)

	if config.Kafka.Enable {
		logger.Debug("kafka exporter config",
			zap.Strings("brokers", config.Kafka.Brokers),
			zap.String("topic", config.Kafka.Topic),
			zap.Duration("timeout", config.Kafka.Timeout),
			zap.Int("retry", config.Kafka.Retry),
			zap.Duration("retryBackoff", config.Kafka.RetryBackoff),
			zap.Int("flushMsg", config.Kafka.FlushMsg),
			zap.Duration("flushTime", config.Kafka.FlushTime),
			zap.Int("flushByte", config.Kafka.FlushByte),
		)
	}

	if config.ElasticSearch.Enable {
		logger.Debug("elasticsearch exporter config",
			zap.Strings("addresses", config.ElasticSearch.Addresses),
			zap.String("user", config.ElasticSearch.User),
			zap.String("pass", config.ElasticSearch.Pass),
			zap.String("index", config.ElasticSearch.Index),
		)
	}

	if config.Volume.Enable {
		logger.Debug("volume config",
			zap.String("fileName", config.Volume.FileName),
			zap.String("filePath", config.Volume.FilePath),
			zap.Int("maxFileSize", config.Volume.MaxFileSize),
			zap.Int("maxFileCount", config.Volume.MaxFileCount),
		)
	}
}

// checkConfig 설정 파일의 필수 항목이 있는지 확인
func checkConfig(config *Config) error {
	if !config.Kafka.Enable && !config.ElasticSearch.Enable && !config.Volume.Enable {
		return fmt.Errorf("at least one exporter is required")
	}

	if config.Kafka.Enable {
		if len(config.Kafka.Brokers) == 0 {
			return fmt.Errorf("kafka brokers is required")
		}
		if config.Kafka.Topic == "" {
			return fmt.Errorf("kafka topic is required")
		}
	}

	if config.ElasticSearch.Enable {
		if len(config.ElasticSearch.Addresses) == 0 {
			return fmt.Errorf("elasticsearch addresses is required")
		}
		if config.ElasticSearch.User == "" {
			return fmt.Errorf("elasticsearch user is required")
		}

		if config.ElasticSearch.Pass == "" {
			return fmt.Errorf("elasticsearch pass is required")
		}

		if config.ElasticSearch.Index == "" {
			return fmt.Errorf("elasticsearch index is required")
		}
	}

	if config.Volume.Enable {
		if config.Volume.FileName == "" {
			return fmt.Errorf("volume fileName is required")
		}

		if config.Volume.FilePath == "" {
			return fmt.Errorf("volume filePath is required")
		}
	}

	return nil
}
