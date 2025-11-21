package exporter

import (
	"context"
	"sync"
)

type Component interface {
	Start(ctx context.Context, wg *sync.WaitGroup) error
}

type Exporter interface {
	Component
	Write(data []byte)
}
