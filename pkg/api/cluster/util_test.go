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

package cluster

import (
	"github.com/huanwei/rocketmq-operator/pkg/apis/rocketmq/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var status = v1alpha1.BrokerClusterStatus{
	TypeMeta:   metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{},
	Conditions: []v1alpha1.BrokerClusterCondition{
		{
			Type:               "TestType",
			Status:             "True",
			LastTransitionTime: metav1.Time{},
			Reason:             "TestReason",
			Message:            "TestMessage",
		},
	},
}

func TestGetClusterCondition(t *testing.T) {
	testStatus := status
	tests := []struct {
		status        *v1alpha1.BrokerClusterStatus
		conditionType v1alpha1.BrokerClusterConditionType
		index         int
		condition     *v1alpha1.BrokerClusterCondition
	}{
		{
			status:        nil,
			conditionType: "",
			index:         -1,
			condition:     nil,
		},
		{
			status:        &testStatus,
			conditionType: "TestType",
			index:         0,
			condition:     &testStatus.Conditions[0],
		},
		{
			status:        &testStatus,
			conditionType: "TestType2",
			index:         -1,
			condition:     nil,
		},
	}
	for i, test := range tests {
		index, condition := GetClusterCondition(test.status, test.conditionType)
		if index != test.index || condition != test.condition {
			t.Errorf("GetClusterCondition Error.\n"+
				"Test number:%d\n"+
				"Expected index:%d,condition:%+v\n"+
				"Get index:%d, condition:%+v\n",
				i+1, test.index, test.condition, index, condition)
		}
	}
}

func TestUpdateClusterCondition(t *testing.T) {
	testStatus := status
	tests := []struct {
		status    *v1alpha1.BrokerClusterStatus
		condition *v1alpha1.BrokerClusterCondition
		result    bool
	}{
		{
			status:    &testStatus,
			condition: &testStatus.Conditions[0],
			result:    false,
		},
		{
			status: &testStatus,
			condition: &v1alpha1.BrokerClusterCondition{
				Type:               "TestType",
				Status:             "False",
				LastTransitionTime: metav1.Time{},
				Reason:             "TestReason2",
				Message:            "TestMessage2",
			},
			result: true,
		},
		{
			status: &testStatus,
			condition: &v1alpha1.BrokerClusterCondition{
				Type:               "TestType2",
				Status:             "True",
				LastTransitionTime: metav1.Time{},
				Reason:             "TestReason2",
				Message:            "TestMessage2",
			},
			result: true,
		},
	}
	for i, test := range tests {
		result := UpdateClusterCondition(test.status, test.condition)
		if result != test.result {
			t.Errorf("UpdateClusterCondition Error.\n"+
				"Test number:%d\n"+
				"Expected result:%t\n"+
				"Get result:%t\n",
				i+1, test.result, result)
		}
	}
}

func TestGetClusterReadyCondition(t *testing.T) {
	testStatus := status
	readyStatus := status
	readyStatus.Conditions = append(readyStatus.Conditions,
		v1alpha1.BrokerClusterCondition{
			Type:               v1alpha1.BrokerClusterReady,
			Status:             "True",
			LastTransitionTime: metav1.Time{},
			Reason:             "TestReason",
			Message:            "TestMessage",
		},
	)
	tests := []struct {
		status    v1alpha1.BrokerClusterStatus
		condition *v1alpha1.BrokerClusterCondition
	}{
		{
			status:    testStatus,
			condition: nil,
		},
		{
			status:    readyStatus,
			condition: &readyStatus.Conditions[1],
		},
	}
	for i, test := range tests {
		condition := GetClusterReadyCondition(test.status)
		if condition != test.condition {
			t.Errorf("GetClusterReadyCondition Error.\n"+
				"Test number:%d\n"+
				"Expected result:%+v\n"+
				"Get result:%+v\n",
				i+1, test.condition, condition)
		}
	}
}

func TestIsClusterReadyConditionTrue(t *testing.T) {
	trueStatus := status
	trueStatus.Conditions = append(trueStatus.Conditions,
		v1alpha1.BrokerClusterCondition{
			Type:               v1alpha1.BrokerClusterReady,
			Status:             "True",
			LastTransitionTime: metav1.Time{},
			Reason:             "TestReason",
			Message:            "TestMessage",
		},
	)
	falseStatus := status
	falseStatus.Conditions = append(falseStatus.Conditions,
		v1alpha1.BrokerClusterCondition{
			Type:               v1alpha1.BrokerClusterReady,
			Status:             "False",
			LastTransitionTime: metav1.Time{},
			Reason:             "TestReason",
			Message:            "TestMessage",
		},
	)
	tests := []struct {
		status v1alpha1.BrokerClusterStatus
		result bool
	}{
		{
			status: trueStatus,
			result: true,
		},
		{
			status: falseStatus,
			result: false,
		},
	}
	for i, test := range tests {
		result := IsClusterReadyConditionTrue(test.status)
		if result != test.result {
			t.Errorf("IsClusterReadyConditionTrue Error.\n"+
				"Test number:%d\n"+
				"Expected result:%+v\n"+
				"Get result:%+v\n",
				i+1, test.result, result)
		}
	}
}
