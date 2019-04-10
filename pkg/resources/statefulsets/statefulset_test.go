/*
Copyright 2017 The Kubernetes Authors.

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
	"encoding/json"
	"github.com/huanwei/rocketmq-operator/pkg/apis/rocketmq/v1alpha1"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
	"unsafe"
)

func TestNewStatefulSet(t *testing.T) {
	tests := []struct {
		cluster *v1alpha1.BrokerCluster
		index   int
		result  *apps.StatefulSet
	}{
		{
			cluster: &v1alpha1.BrokerCluster{},
			index:   0,
			result:  &apps.StatefulSet{},
		},
		{
			cluster: &v1alpha1.BrokerCluster{},
			index:   1,
			result:  &apps.StatefulSet{},
		},
	}
	for i, test := range tests {
		result := NewStatefulSet(test.cluster, test.index)
		out, _ := json.MarshalIndent(result, "", "     ")
		str := (*string)(unsafe.Pointer(&out))
		t.Logf("Test number:%d\n"+
			"Get result:\n%s\n",
			i+1, *str)
	}
}

func TestParseQuantity(t *testing.T) {
	q, _ := resource.ParseQuantity("2Gi")
	t.Logf("Get result: %v", q)
}
