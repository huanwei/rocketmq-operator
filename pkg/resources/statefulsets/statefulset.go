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
	"github.com/huanwei/rocketmq-operator/pkg/apis/rocketmq/v1alpha1"
	"github.com/huanwei/rocketmq-operator/pkg/constants"
	apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

	var logQuantity,storeQuantity  resource.Quantity
	var err error
	logQuantity, err = resource.ParseQuantity("5Gi")
	if err != nil {
		return nil
	}
	storeQuantity, err = resource.ParseQuantity("5Gi")
	if err != nil {
		return nil
	}

	storageClassNmae := cluster.Spec.StorageClassName
	if storageClassNmae == "" {
		storageClassNmae = constants.DefaultStorageClassName
	}
	volumeClaimTemplates := []v1.PersistentVolumeClaim{
		nfsPersistentVolumeClaim(storageClassNmae, logQuantity, "brokerlogs"),
		nfsPersistentVolumeClaim(storageClassNmae, storeQuantity, "brokerstore"),
	}

	ssReplicas := int32(cluster.Spec.MembersPerGroup)
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
		Spec: apps.StatefulSetSpec{
			Replicas: &ssReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: v1.PodSpec{
					//ServiceAccountName: "rocketmq-operator",
					NodeSelector: cluster.Spec.NodeSelector,
					//Affinity:     cluster.Spec.Affinity,
					Containers: containers,
				},
			},
			ServiceName: fmt.Sprintf(cluster.Name+`-svc-%d`, index),
			VolumeClaimTemplates: volumeClaimTemplates,
		},
	}
	return ss

}

func brokerContainer(cluster *v1alpha1.BrokerCluster, index int) v1.Container {
	return v1.Container{
		Name:            "broker",
		ImagePullPolicy: "Always",
		Image:           cluster.Spec.BrokerImage,
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
				MountPath: "/root/logs",
			},
			{
				Name:      "brokerstore",
				MountPath: "/root/store",
			},
		},
	}

}

func nfsPersistentVolumeClaim(storageClassName string, quantity resource.Quantity, name string) v1.PersistentVolumeClaim{
	annotations := map[string]string{
		constants.BrokerVolumeStorageClass: storageClassName,
	}
	return v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: annotations,
		},
		Spec: v1.PersistentVolumeClaimSpec {
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: quantity,
				},
			},
		},
	}
}
