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
	"github.com/huanwei/rocketmq-operator/pkg/test/mock"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"testing"
)

type updateCallback func(cluster *v1alpha1.BrokerCluster) (*v1alpha1.BrokerCluster, error)
type getCallback func(name string) (*v1alpha1.BrokerCluster, error)

var updateError = errors.New("update fail")
var getError = errors.New("get fail")
var brokerCluster = &v1alpha1.BrokerCluster{}

func initTest() (*clusterUpdater, *mock.ClientSet, *mock.BrokerClusterLister) {
	cu := clusterUpdater{}
	mockClient := mock.NewMockClient()
	mockLister := mock.NewMockBrokerClusterLister()
	cu.client = mockClient
	cu.lister = mockLister

	return &cu, mockClient, mockLister
}

func updateSuccess(cluster *v1alpha1.BrokerCluster) (*v1alpha1.BrokerCluster, error) {
	return nil, nil
}

func updateFail(cluster *v1alpha1.BrokerCluster) (*v1alpha1.BrokerCluster, error) {
	return nil, updateError
}

func getSuccess(name string) (*v1alpha1.BrokerCluster, error) {
	return nil, nil
}

func getFail(name string) (*v1alpha1.BrokerCluster, error) {
	return nil, getError
}

func TestClusterUpdater_UpdateClusterStatus(t *testing.T) {
	cu, client, lister := initTest()
	tests := []struct {
		update updateCallback
		get    getCallback
		err    error
	}{
		{
			update: updateFail,
			get:    getSuccess,
			err:    updateError,
		},
		{
			update: updateFail,
			get:    getFail,
			err:    getError,
		},
		{
			update: updateSuccess,
			get:    nil,
			err:    nil,
		},
	}
	for i, test := range tests {
		client.ROCKETMQClient.BrokerClustersRef.UpdateCallback = test.update
		lister.NsLister.GetCallback = test.get
		err := cu.UpdateClusterStatus(brokerCluster, &v1alpha1.BrokerClusterStatus{})
		if err != test.err {
			t.Errorf("UpdateClusterStatus Error.\n"+
				"Test number:%d\n"+
				"Expected result:%+v\n"+
				"Get result:%+v\n",
				i+1, test.err, err)
		}
	}
}

func TestClusterUpdater_UpdateClusterLabels(t *testing.T) {
	cu, client, lister := initTest()
	tests := []struct {
		update updateCallback
		get    getCallback
		err    error
	}{
		{
			update: updateFail,
			get:    getSuccess,
			err:    updateError,
		},
		{
			update: updateFail,
			get:    getFail,
			err:    getError,
		},
		{
			update: updateSuccess,
			get:    nil,
			err:    nil,
		},
	}
	for i, test := range tests {
		client.ROCKETMQClient.BrokerClustersRef.UpdateCallback = test.update
		lister.NsLister.GetCallback = test.get
		err := cu.UpdateClusterLabels(brokerCluster, labels.Set{})
		if err != test.err {
			t.Errorf("UpdateClusterLabels Error.\n"+
				"Test number:%d\n"+
				"Expected result:%+v\n"+
				"Get result:%+v\n",
				i+1, test.err, err)
		}
	}
}
