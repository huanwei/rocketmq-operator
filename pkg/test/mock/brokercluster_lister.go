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

package mock

import (
	"github.com/huanwei/rocketmq-operator/pkg/apis/rocketmq/v1alpha1"
	listersv1alpha1 "github.com/huanwei/rocketmq-operator/pkg/generated/listers/rocketmq/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
)

type BrokerClusterLister struct {
	ListCallback func(selector labels.Selector) (ret []*v1alpha1.BrokerCluster, err error)
	NsLister     *BrokerClusterNamespaceLister
}

type BrokerClusterNamespaceLister struct {
	ListCallback func(selector labels.Selector) (ret []*v1alpha1.BrokerCluster, err error)
	GetCallback  func(name string) (*v1alpha1.BrokerCluster, error)
}

func NewMockBrokerClusterLister() *BrokerClusterLister {
	return &BrokerClusterLister{
		ListCallback: nil,
		NsLister:     NewMockBrokerClusterNamespaceLister(),
	}
}

func (s *BrokerClusterLister) List(selector labels.Selector) (ret []*v1alpha1.BrokerCluster, err error) {
	return s.ListCallback(selector)
}

func (s *BrokerClusterLister) BrokerClusters(namespace string) listersv1alpha1.BrokerClusterNamespaceLister {
	return s.NsLister
}

func NewMockBrokerClusterNamespaceLister() *BrokerClusterNamespaceLister {
	return &BrokerClusterNamespaceLister{
		ListCallback: nil,
		GetCallback:  nil,
	}
}

func (s BrokerClusterNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.BrokerCluster, err error) {
	return s.ListCallback(selector)
}

func (s BrokerClusterNamespaceLister) Get(name string) (*v1alpha1.BrokerCluster, error) {
	return s.GetCallback(name)
}
