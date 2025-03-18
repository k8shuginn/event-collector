package exporter

import "context"

type Component interface {
	Start(ctx context.Context) error
}

type Exporter interface {
	Component
	Write(data []byte)
}
