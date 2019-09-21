package leaderelection

import (
	"fmt"

	coordinationv1 "k8s.io/api/coordination/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type LockType int

const (
	UndefinedResourceLock LockType = iota
	ConfigMapsResourceLock
	LeasesResourceLock
	EndpointsResourceLock
)

func (l LockType) Name() string {
	switch l {
	case EndpointsResourceLock:
		return resourcelock.EndpointsResourceLock
	case ConfigMapsResourceLock:
		return resourcelock.ConfigMapsResourceLock
	case LeasesResourceLock:
		return resourcelock.LeasesResourceLock
	default:
		return ""
	}
}

// GetPreferredLockType chooses the Lease lock if `lease.coordination.k8s.io` is available.
// Otherwise, the ConfigMap resource lock is used.
func GetPreferredLockType(mapper meta.RESTMapper) (LockType, error) {
	// check if new leader election api is available
	_, err := mapper.RESTMapping(schema.GroupKind{
		Kind:  "Lease",
		Group: coordinationv1.GroupName,
	})
	switch {
	case err == nil:
		return LeasesResourceLock, nil
	case meta.IsNoMatchError(err):
		return ConfigMapsResourceLock, nil
	default:
		return UndefinedResourceLock, fmt.Errorf("unable to retrieve supported server groups: %v", err)
	}
}
