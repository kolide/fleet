// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import (
	"time"

	"github.com/kolide/fleet/server/kolide"
)

var _ kolide.TargetStore = (*TargetStore)(nil)

type CountHostsInTargetsFunc func(hostIDs, labelIDs []uint, now time.Time) (kolide.TargetMetrics, error)
type HostIDsInTargetsFunc func(hostIDs, labelIDs []uint) ([]uint, error)

type TargetStore struct {
	CountHostsInTargetsFunc        CountHostsInTargetsFunc
	CountHostsInTargetsFuncInvoked bool
	HostIDsInTargetsFunc           HostIDsInTargetsFunc
	HostIDsInTargetsFuncInvoked    bool
}

func (s *TargetStore) CountHostsInTargets(hostIDs, labelIDs []uint, now time.Time) (kolide.TargetMetrics, error) {
	s.CountHostsInTargetsFuncInvoked = true
	return s.CountHostsInTargetsFunc(hostIDs, labelIDs, now)
}

func (s *TargetStore) HostIDsInTargets(hostIDs, labelIDs []uint) ([]uint, error) {
	s.HostIDsInTargetsFuncInvoked = true
	return s.HostIDsInTargetsFunc(hostIDs, labelIDs)
}
