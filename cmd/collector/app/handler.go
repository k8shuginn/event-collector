package app

import (
	"github.com/k8shuginn/event-collector/exporter"
	"github.com/k8shuginn/event-collector/pkg/kube"
	"github.com/k8shuginn/event-collector/pkg/logger"
	"go.uber.org/zap"
	v1 "k8s.io/api/events/v1"
)

type Handler struct {
	exporters []exporter.Exporter
}

// OnAdd event handler
func (h *Handler) OnAdd(obj interface{}, _ bool) {
	object := obj.(*v1.Event)
	event, err := kube.ConvertBytes(object)
	if err != nil {
		logger.Error("failed to convert OnAdd event object", zap.Error(err))
		return
	}

	logger.Debug("kubernetes event OnAdd",
		zap.String("event", object.Name),
		zap.String("namespace", object.Namespace),
		zap.String("kind", object.Regarding.Kind),
	)

	for _, e := range h.exporters {
		e.Write(event)
	}
}

// OnUpdate event handler
func (h *Handler) OnUpdate(oldObj, newObj interface{}) {
	object := newObj.(*v1.Event)
	event, err := kube.ConvertBytes(object)
	if err != nil {
		logger.Error("failed to convert OnUpdate event object", zap.Error(err))
		return
	}

	logger.Debug("kubernetes event OnUpdate",
		zap.String("event", object.Name),
		zap.String("namespace", object.Namespace),
		zap.String("kind", object.Regarding.Kind),
	)

	for _, e := range h.exporters {
		e.Write(event)
	}
}

// OnDelete event handler
func (h *Handler) OnDelete(obj interface{}) {
	object := obj.(*v1.Event)
	event, err := kube.ConvertBytes(object)
	if err != nil {
		logger.Error("failed to convert OnDelete event object", zap.Error(err))
		return
	}

	logger.Debug("kubernetes event OnDelete",
		zap.String("event", object.Name),
		zap.String("namespace", object.Namespace),
		zap.String("kind", object.Regarding.Kind),
	)

	for _, e := range h.exporters {
		e.Write(event)
	}
}
