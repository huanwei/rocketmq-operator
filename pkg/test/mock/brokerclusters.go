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
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type BrokerClusters struct {
	CreateCallback           func(*v1alpha1.BrokerCluster) (*v1alpha1.BrokerCluster, error)
	UpdateCallback           func(*v1alpha1.BrokerCluster) (*v1alpha1.BrokerCluster, error)
	DeleteCallback           func(name string, options *v1.DeleteOptions) error
	DeleteCollectionCallback func(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	GetCallback              func(name string, options v1.GetOptions) (*v1alpha1.BrokerCluster, error)
	ListCallback             func(opts v1.ListOptions) (*v1alpha1.BrokerClusterList, error)
	WatchCallback            func(opts v1.ListOptions) (watch.Interface, error)
	PatchCallback            func(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BrokerCluster, err error)
}

func (bc *BrokerClusters) Create(c *v1alpha1.BrokerCluster) (*v1alpha1.BrokerCluster, error) {
	return bc.CreateCallback(c)
}

func (bc *BrokerClusters) Update(c *v1alpha1.BrokerCluster) (*v1alpha1.BrokerCluster, error) {
	return bc.UpdateCallback(c)
}

func (bc *BrokerClusters) Delete(name string, options *v1.DeleteOptions) error {
	return bc.DeleteCallback(name, options)
}

func (bc *BrokerClusters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return bc.DeleteCollectionCallback(options, listOptions)
}

func (bc *BrokerClusters) Get(name string, options v1.GetOptions) (*v1alpha1.BrokerCluster, error) {
	return bc.GetCallback(name, options)
}

func (bc *BrokerClusters) List(opts v1.ListOptions) (*v1alpha1.BrokerClusterList, error) {
	return bc.ListCallback(opts)
}

func (bc *BrokerClusters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return bc.WatchCallback(opts)
}

func (bc *BrokerClusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.BrokerCluster, err error) {
	return bc.PatchCallback(name, pt, data, "")
}
