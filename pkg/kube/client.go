package kube

import (
	"fmt"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type ifm struct {
	factory  informers.SharedInformerFactory
	informer cache.SharedIndexInformer
}

type Client struct {
	cs *kubernetes.Clientset

	closeCh chan struct{}
	ifms    map[string]ifm
}

func NewClient(eh cache.ResourceEventHandler, opts ...Option) (*Client, error) {
	c := fromOptions(opts...)
	client := &Client{
		ifms: make(map[string]ifm),
	}

	// clientConfig 설정
	var clientConfig *rest.Config
	var err error
	if c.kubeConfig != "" {
		clientConfig, err = clientcmd.BuildConfigFromFlags("", c.kubeConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubernetes config: %v", err)
		}
	} else {
		clientConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to build kubernetes in-cluster config: %v", err)
		}
	}

	// clientset 생성
	client.cs, err = kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %v", err)
	}

	if len(c.namespaces) > 0 {
		// namespace별로 informer 생성
		for ns := range c.namespaces {
			factory := informers.NewSharedInformerFactoryWithOptions(client.cs, c.resyncTime, informers.WithNamespace(ns))
			informer := factory.Events().V1().Events().Informer()
			informer.AddEventHandler(eh)
			client.ifms[ns] = ifm{
				factory:  factory,
				informer: informer,
			}
		}
	} else {
		// 모든 namespace를 수집하는 경우
		factory := informers.NewSharedInformerFactory(client.cs, c.resyncTime)
		informer := factory.Events().V1().Events().Informer()
		informer.AddEventHandler(eh)
		client.ifms["*"] = ifm{
			factory:  factory,
			informer: informer,
		}
	}

	return client, nil
}

// Run client 실행
func (c *Client) Run() {
	c.closeCh = make(chan struct{})

	for _, ifm := range c.ifms {
		go ifm.factory.Start(c.closeCh)
	}
}

// Close client 종료
func (c *Client) Close() {
	close(c.closeCh)
}
