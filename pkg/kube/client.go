package kube

import (
	"fmt"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	cs       *kubernetes.Clientset
	iFactory informers.SharedInformerFactory
	closeCh  chan struct{}

	sii cache.SharedIndexInformer              // SharedIndexInformer 객체
	reg cache.ResourceEventHandlerRegistration // event handler 등록 정보
}

func NewClient(eh cache.ResourceEventHandler, opts ...Option) (*Client, error) {
	c := fromOptions(opts...)
	client := &Client{}

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

	// SharedInformerFactory 생성
	client.iFactory = informers.NewSharedInformerFactory(client.cs, c.resyncTime)

	// SharedIndexInformer, ResourceEventHandler 생성
	client.sii = client.iFactory.Events().V1().Events().Informer()
	client.reg, err = client.sii.AddEventHandler(eh)
	if err != nil {
		return nil, fmt.Errorf("failed to add kubernetes event handler: %v", err)
	}

	return client, nil
}

// Run client 실행
func (c *Client) Run() {
	c.closeCh = make(chan struct{})
	go c.iFactory.Start(c.closeCh)
}

// Close client 종료
func (c *Client) Close() {
	_ = c.sii.RemoveEventHandler(c.reg)
	close(c.closeCh)
}
