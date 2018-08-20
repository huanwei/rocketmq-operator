package v1alpha1

import corev1 "k8s.io/api/core/v1"
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const MinimumRocketMQVersion = "4.2.0"

type BrokerClusterSpec struct {
	Version string `json:"version"`
	Namesrvs string `json:"namesrvs"`
	Members int32 `json:"members, omitempty"`
	BaseBrokerID uint32 `json:"baseBrokerId, omitempty"`
	ClusterMode bool `json:"clusterMode, omitempty`
	NodeSelector map[string]string `json:"nodeSelector, omitempty"`
	Affinity *corev1.Affinity `json:"affinity, omitempty"`
	VolumeClaimTemplate *corev1.PersistentVolumeClaim `json:"volumeClaimTemplate, omitempty"`
	Config *corev1.LocalObjectReference `json:"config,omitempty"`
}

type BrokerClusterConditonType string

const BrokerClusterReady  BrokerClusterConditonType = "Ready"

type BrokerClusterCondition struct {
	Type BrokerClusterConditonType
	Status corev1.ConditionStatus
	LastTransitionTime metav1.Time
	Reason string
	Message string
}

type BrokerClusterStatus struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Conditions []BrokerClusterCondition
}

// +genclient
// +genclient:noStatus
// +resourceName=brokerclusters
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BrokerCluster struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec BrokerClusterSpec `json:"spec"`
	Status BrokerClusterStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BrokerClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items []BrokerCluster `json:"items"`
}