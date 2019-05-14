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

package statefulsets

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/huanwei/rocketmq-operator/pkg/apis/rocketmq/v1alpha1"
	"github.com/huanwei/rocketmq-operator/pkg/constants"
	apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strconv"
)

func NewStatefulSet(cluster *v1alpha1.BrokerCluster, index int) *apps.StatefulSet {
	containers := []v1.Container{
		brokerContainer(cluster, index),
	}

	brokerRole := constants.BrokerRoleSlave
	if index == 0 {
		brokerRole = constants.BrokerRoleMaster
	}
	podLabels := map[string]string{
		constants.BrokerClusterLabel: fmt.Sprintf(cluster.Name+`-%d`, index),
		constants.BrokerRoleLabel:    brokerRole,
	}

	podSpec := v1.PodSpec{
		NodeSelector: cluster.Spec.NodeSelector,
		Containers:   containers,
	}

	if cluster.Spec.EmptyDir {
		podSpec = v1.PodSpec{
			NodeSelector: cluster.Spec.NodeSelector,
			Containers:   containers,
			Volumes: []v1.Volume{
				{
					Name: "brokerlogs",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: "brokerstore",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
			},
		}
	}

	ssReplicas := int32(cluster.Spec.MembersPerGroup)

	statefulsetSpec := apps.StatefulSetSpec{
		Replicas: &ssReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: podLabels,
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: podLabels,
			},
			Spec: podSpec,
		},
		UpdateStrategy: apps.StatefulSetUpdateStrategy{
			Type: apps.RollingUpdateStatefulSetStrategyType,
		},
		ServiceName: fmt.Sprintf(cluster.Name+`-svc-%d`, index),
	}

	var logQuantity, storeQuantity resource.Quantity
	var err error
	logQuantity = defaultStorageQuantity()
	storeQuantity = defaultStorageQuantity()

	if cluster.Spec.ContainerSpec.Requests != nil {
		logSize := cluster.Spec.ContainerSpec.Requests.LogStorage
		if logSize != "" {
			logQuantity, err = resource.ParseQuantity(logSize)
			if err != nil {
				glog.Errorf("failed to parse log size %s to quantity: %v", logSize, err)
				return nil
			}
		}
		storeSize := cluster.Spec.ContainerSpec.Requests.StoreStorage
		if storeSize != "" {
			storeQuantity, err = resource.ParseQuantity(storeSize)
			if err != nil {
				glog.Errorf("failed to parse store size %s to quantity: %v", storeSize, err)
				return nil
			}
		}
	}

	storageClassNmae := cluster.Spec.StorageClassName
	if storageClassNmae == "" {
		storageClassNmae = constants.DefaultStorageClassName
	}
	volumeClaimTemplates := []v1.PersistentVolumeClaim{
		nfsPersistentVolumeClaim(storageClassNmae, logQuantity, "brokerlogs"),
		nfsPersistentVolumeClaim(storageClassNmae, storeQuantity, "brokerstore"),
	}

	if !cluster.Spec.EmptyDir && storageClassNmae != "" {
		statefulsetSpec = apps.StatefulSetSpec{
			Replicas: &ssReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: podSpec,
			},
			UpdateStrategy: apps.StatefulSetUpdateStrategy{
				Type: apps.RollingUpdateStatefulSetStrategyType,
			},
			ServiceName:          fmt.Sprintf(cluster.Name+`-svc-%d`, index),
			VolumeClaimTemplates: volumeClaimTemplates,
		}
	}

	ss := &apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      fmt.Sprintf(cluster.Name+`-%d`, index),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cluster, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    v1alpha1.ClusterCRDResourceKind,
				}),
			},
			Labels: podLabels,
		},
		Spec: statefulsetSpec,
	}
	return ss

}

func brokerContainer(cluster *v1alpha1.BrokerCluster, index int) v1.Container {
	return v1.Container{
		Name:            "broker",
		ImagePullPolicy: "Always",
		Image:           cluster.Spec.ContainerSpec.BrokerImage,
		Ports: []v1.ContainerPort{
			{
				ContainerPort: 10909,
			},
			{
				ContainerPort: 10911,
			},
		},
		Env: []v1.EnvVar{
			{
				Name:  "DELETE_WHEN",
				Value: cluster.Spec.Properties["DELETE_WHEN"],
			},
			{
				Name:  "FILE_RESERVED_TIME",
				Value: cluster.Spec.Properties["FILE_RESERVED_TIME"],
			},
			{
				Name:  "ALL_MASTER",
				Value: strconv.FormatBool(cluster.Spec.AllMaster),
			},
			{
				Name:  "BROKER_NAME",
				Value: fmt.Sprintf(`broker-%d`, index),
			},
			{
				Name:  "REPLICATION_MODE",
				Value: cluster.Spec.ReplicationMode,
			},
			{
				Name:  "FLUSH_DISK_TYPE",
				Value: cluster.Spec.Properties["FLUSH_DISK_TYPE"],
			},
			{
				Name:  "NAMESRV_ADDRESS",
				Value: cluster.Spec.NameServers,
			},
		},
		Command: []string{"./brokerStart.sh"},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "brokerlogs",
				MountPath: "/home/rocketmq/logs",
			},
			{
				Name:      "brokerstore",
				MountPath: "/home/rocketmq/store",
			},
		},
		Resources: resourceRequirement(cluster.Spec.ContainerSpec, defaultRequests()),
	}

}

func resourceRequirement(spec v1alpha1.ContainerSpec, defaultRequests ...v1.ResourceRequirements) v1.ResourceRequirements {
	rr := v1.ResourceRequirements{}
	if len(defaultRequests) > 0 {
		defaultRequest := defaultRequests[0]
		rr.Requests = make(map[v1.ResourceName]resource.Quantity)
		rr.Requests[v1.ResourceCPU] = defaultRequest.Requests[v1.ResourceCPU]
		rr.Requests[v1.ResourceMemory] = defaultRequest.Requests[v1.ResourceMemory]
		rr.Limits = make(map[v1.ResourceName]resource.Quantity)
		rr.Limits[v1.ResourceCPU] = defaultRequest.Limits[v1.ResourceCPU]
		rr.Limits[v1.ResourceMemory] = defaultRequest.Limits[v1.ResourceMemory]
	}
	if spec.Requests != nil {
		if rr.Requests == nil {
			rr.Requests = make(map[v1.ResourceName]resource.Quantity)
		}
		if spec.Requests.CPU != "" {
			if q, err := resource.ParseQuantity(spec.Requests.CPU); err != nil {
				glog.Errorf("failed to parse CPU resource %s to quantity: %v", spec.Requests.CPU, err)
			} else {
				rr.Requests[v1.ResourceCPU] = q
			}
		}
		if spec.Requests.Memory != "" {
			if q, err := resource.ParseQuantity(spec.Requests.Memory); err != nil {
				glog.Errorf("failed to parse memory resource %s to quantity: %v", spec.Requests.Memory, err)
			} else {
				rr.Requests[v1.ResourceMemory] = q
			}
		}
	}
	if spec.Limits != nil {
		if rr.Limits == nil {
			rr.Limits = make(map[v1.ResourceName]resource.Quantity)
		}
		if spec.Limits.CPU != "" {
			if q, err := resource.ParseQuantity(spec.Limits.CPU); err != nil {
				glog.Errorf("failed to parse CPU resource %s to quantity: %v", spec.Limits.CPU, err)
			} else {
				rr.Limits[v1.ResourceCPU] = q
			}
		}
		if spec.Limits.Memory != "" {
			if q, err := resource.ParseQuantity(spec.Limits.Memory); err != nil {
				glog.Errorf("failed to parse memory resource %s to quantity: %v", spec.Limits.Memory, err)
			} else {
				rr.Limits[v1.ResourceMemory] = q
			}
		}
	}
	return rr

}

func defaultRequests() v1.ResourceRequirements {
	rr := v1.ResourceRequirements{}
	rr.Requests = make(map[v1.ResourceName]resource.Quantity)
	rr.Limits = make(map[v1.ResourceName]resource.Quantity)
	rr.Requests[v1.ResourceCPU] = resource.MustParse("1000m")
	rr.Requests[v1.ResourceMemory] = resource.MustParse("1000Mi")
	rr.Limits[v1.ResourceCPU] = resource.MustParse("2000m")
	rr.Limits[v1.ResourceMemory] = resource.MustParse("2000Mi")
	return rr
}

func defaultStorageQuantity() resource.Quantity {
	var q resource.Quantity
	q, _ = resource.ParseQuantity("5Gi")
	return q
}

func nfsPersistentVolumeClaim(storageClassName string, quantity resource.Quantity, name string) v1.PersistentVolumeClaim {
	annotations := map[string]string{
		constants.BrokerVolumeStorageClass: storageClassName,
	}
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: annotations,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: quantity,
				},
			},
		},
	}
}
