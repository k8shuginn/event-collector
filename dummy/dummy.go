package dummy

import (
	"time"

	"github.com/k8shuginn/event-collector/pkg/kube"
)

func MakeDummy(reVersion string) *kube.Event {
	dummy := kube.Event{}
	dummy.Metadata.Name = "test"
	dummy.Metadata.Namespace = "gh-runner"
	dummy.Metadata.UID = "test1234"
	dummy.Metadata.ResourceVersion = reVersion
	dummy.Metadata.CreationTimestamp = time.Now()

	dummy.EventTime = time.Now()
	dummy.RetportingController = "test"
	dummy.Reason = "test"

	dummy.Regarding.Kind = "test"
	dummy.Regarding.Namespace = "test"
	dummy.Regarding.Name = "test"
	dummy.Regarding.UID = "test1234"
	dummy.Regarding.ApiVersion = "test.io"
	dummy.Regarding.ResourceVersion = reVersion

	dummy.Note = "test create"
	dummy.Type = "Normal"

	dummy.DeprecatedFirstTimestamp = time.Now()
	dummy.DeprecatedLastTimestamp = time.Now()
	dummy.DeprecatedCount = 1

	return &dummy
}
