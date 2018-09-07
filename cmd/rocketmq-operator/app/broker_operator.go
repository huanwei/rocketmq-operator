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

package app

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/glog"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	cluster "github.com/huanwei/rocketmq-operator/pkg/controllers/cluster"
	clientset "github.com/huanwei/rocketmq-operator/pkg/generated/clientset/versioned"
	informers "github.com/huanwei/rocketmq-operator/pkg/generated/informers/externalversions"
	operatoropts "github.com/huanwei/rocketmq-operator/pkg/options"
	signals "github.com/huanwei/rocketmq-operator/pkg/util/signals"
)

// resyncPeriod computes the time interval a shared informer waits before
// resyncing with the api server.
func resyncPeriod(s *operatoropts.OperatorOpts) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(s.MinResyncPeriod.Nanoseconds()) * factor)
	}
}

// Run starts the operator controllers. This should never exit.
func Run(s *operatoropts.OperatorOpts) error {
	kubeconfig, err := clientcmd.BuildConfigFromFlags(s.Master, s.KubeConfig)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// Set up signals so we handle the first shutdown signal gracefully.
	signals.SetupSignalHandler(cancelFunc)

	kubeClient := kubernetes.NewForConfigOrDie(kubeconfig)
	opClient := clientset.NewForConfigOrDie(kubeconfig)

	// Shared informers (non namespace specific).
	kubeInformerFactory := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, resyncPeriod(s)(), s.Namespace, nil)
	operatorInformerFactory := informers.NewFilteredSharedInformerFactory(opClient, resyncPeriod(s)(), s.Namespace, nil)

	var wg sync.WaitGroup

	clusterController := cluster.NewBrokerController(
		*s,
		opClient,
		kubeClient,
		operatorInformerFactory.ROCKETMQ().V1alpha1().BrokerClusters(),
		kubeInformerFactory.Apps().V1().StatefulSets(),
		kubeInformerFactory.Core().V1().Pods(),
		kubeInformerFactory.Core().V1().Services(),
		30*time.Second,
		s.Namespace,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		clusterController.Run(ctx, 5)
	}()

	// Shared informers have to be started after ALL controllers.
	go kubeInformerFactory.Start(ctx.Done())
	go operatorInformerFactory.Start(ctx.Done())

	<-ctx.Done()

	glog.Info("Waiting for all controllers to shut down gracefully")
	wg.Wait()

	return nil
}
