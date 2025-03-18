package kube

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	v1 "k8s.io/api/events/v1"
)

type testHandler struct {
	t *testing.T
}

func (h *testHandler) OnAdd(obj interface{}, isInInitialList bool) {
	object := obj.(*v1.Event)
	event := ConvertEvent(object)
	h.t.Logf("OnAdd: %v", event)
}

func (h *testHandler) OnUpdate(oldObj, newObj interface{}) {
	object := newObj.(*v1.Event)
	event := ConvertEvent(object)
	h.t.Logf("OnUpdate: %v", event)
}

func (h *testHandler) OnDelete(obj interface{}) {
	object := obj.(*v1.Event)
	event := ConvertEvent(object)
	h.t.Logf("OnDelete: %v", event)
}

func TestClient(t *testing.T) {
	handler := &testHandler{
		t: t,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home directory: %v", err)
	}

	client, err := NewClient(
		handler,
		WithKubeConfig(filepath.Join(home, ".kube", "config")),
		WithResyncTime(0),
	)

	client.Run()
	time.Sleep(10 * time.Second)
	client.Close()
}
