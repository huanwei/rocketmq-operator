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

package services

import (
	"fmt"
	"github.com/huanwei/rocketmq-operator/pkg/apis/rocketmq/v1alpha1"
	"github.com/huanwei/rocketmq-operator/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHeadlessService(cluster *v1alpha1.BrokerCluster, index int) *corev1.Service {
	var ports []corev1.ServicePort
	ports = append(ports, corev1.ServicePort{
		Port: 10909,
		Name: "port10909",
	})
	ports = append(ports, corev1.ServicePort{
		Port: 10911,
		Name: "port10911",
	})
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			//Labels:    map[string]string{constants.BrokerClusterLabel: cluster.Name},
			//Name:      cluster.Name,
			Name: fmt.Sprintf(cluster.Name+`-svc-%d`, index),
			Labels: map[string]string{
				constants.BrokerClusterLabel: fmt.Sprintf(cluster.Name+`-%d`, index),
			},
			Namespace: cluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cluster, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    v1alpha1.ClusterCRDResourceKind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				constants.BrokerClusterLabel: fmt.Sprintf(cluster.Name+`-%d`, index),
			},
			ClusterIP: corev1.ClusterIPNone,
		},
	}
	return svc
}
