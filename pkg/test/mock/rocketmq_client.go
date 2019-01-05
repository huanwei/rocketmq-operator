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
	rocketmqv1alpha1 "github.com/huanwei/rocketmq-operator/pkg/generated/clientset/versioned/typed/rocketmq/v1alpha1"
	"k8s.io/client-go/rest"
)

type ROCKETMQV1alpha1Client struct {
	restClient        rest.Interface
	BrokerClustersRef *BrokerClusters
}

func NewMockROCKETMQClient() *ROCKETMQV1alpha1Client {
	c := ROCKETMQV1alpha1Client{}
	c.BrokerClustersRef = &BrokerClusters{
		CreateCallback:           nil,
		UpdateCallback:           nil,
		DeleteCallback:           nil,
		DeleteCollectionCallback: nil,
		GetCallback:              nil,
		ListCallback:             nil,
		WatchCallback:            nil,
		PatchCallback:            nil,
	}
	return &c
}

func (c *ROCKETMQV1alpha1Client) BrokerClusters(namespace string) rocketmqv1alpha1.BrokerClusterInterface {
	return c.BrokerClustersRef
}

func (c *ROCKETMQV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
