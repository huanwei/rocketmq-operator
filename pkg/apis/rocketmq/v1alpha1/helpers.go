/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

const (
	defaultBrokerImage     = "huanwei/rocketmq-broker:4.3.0-operator"
	defaultGroups          = 2
	defaultMembersPerGroup = 2

	defaultDeleteWhen       = "04"
	defaultFileReservedTime = "48"
	defaultReplicationMode  = "SYNC"
	defaultFlushDiskType    = "ASYNC_FLUSH"
)

// EnsureDefaults will ensure that if a user omits and fields in the
// spec that are required, we set some sensible defaults.
// For example a user can choose to omit the version
// and number of members.
func (c *BrokerCluster) EnsureDefaults() *BrokerCluster {
	if c.Spec.BrokerImage == "" {
		c.Spec.BrokerImage = defaultBrokerImage
	}
	if c.Spec.GroupReplica == 0 {
		c.Spec.GroupReplica = defaultGroups
	}

	if c.Spec.MembersPerGroup == 0 {
		c.Spec.MembersPerGroup = defaultMembersPerGroup
	}
	if c.Spec.Properties == nil {
		c.Spec.Properties = map[string]string{}
	}
	if c.Spec.Properties["DELETE_WHEN"] == "" {
		c.Spec.Properties["DELETE_WHEN"] = defaultDeleteWhen
	}
	if c.Spec.Properties["FILE_RESERVED_TIME"] == "" {
		c.Spec.Properties["FILE_RESERVED_TIME"] = defaultFileReservedTime
	}

	if c.Spec.Properties["REPLICATION_MODE"] == "" {
		c.Spec.Properties["REPLICATION_MODE"] = defaultReplicationMode
	}
	if c.Spec.Properties["FLUSH_DISK_TYPE"] == "" {
		c.Spec.Properties["FLUSH_DISK_TYPE"] = defaultFlushDiskType
	}
	return c
}
