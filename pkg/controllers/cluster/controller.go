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
	"context"
	"fmt"
	"time"

	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	wait "k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	scheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
	record "k8s.io/client-go/tools/record"
	workqueue "k8s.io/client-go/util/workqueue"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	clusterutil "github.com/huanwei/rocketmq-operator/pkg/api/cluster"
	v1alpha1 "github.com/huanwei/rocketmq-operator/pkg/apis/rocketmq/v1alpha1"
	constants "github.com/huanwei/rocketmq-operator/pkg/constants"
	controllerutils "github.com/huanwei/rocketmq-operator/pkg/controllers/util"
	clientset "github.com/huanwei/rocketmq-operator/pkg/generated/clientset/versioned"
	opscheme "github.com/huanwei/rocketmq-operator/pkg/generated/clientset/versioned/scheme"
	informersv1alpha1 "github.com/huanwei/rocketmq-operator/pkg/generated/informers/externalversions/rocketmq/v1alpha1"
	listersv1alpha1 "github.com/huanwei/rocketmq-operator/pkg/generated/listers/rocketmq/v1alpha1"

	operatoropts "github.com/huanwei/rocketmq-operator/pkg/options"
	services "github.com/huanwei/rocketmq-operator/pkg/resources/services"
	statefulsets "github.com/huanwei/rocketmq-operator/pkg/resources/statefulsets"
)

const controllerAgentName = "rocketmq-operator"

const (
	// SuccessSynced is used as part of the Event 'reason' when a MySQSL is
	// synced.
	SuccessSynced = "Synced"

	// MessageResourceSynced is the message used for an Event fired when a
	// Cluster is synced successfully
	MessageResourceSynced = "Cluster synced successfully"

	// ErrResourceExists is used as part of the Event 'reason' when a
	// Cluster fails to sync due to a resource of the same name already
	// existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a resource already existing.
	MessageResourceExists = "%s %s/%s already exists and is not managed by Cluster"
)

// The BrokerController watches the Kubernetes API for changes to rocketmq resources
type BrokerController struct {
	// Global Operator configuration options.
	opConfig operatoropts.OperatorOpts

	kubeClient kubernetes.Interface
	opClient   clientset.Interface

	shutdown bool
	queue    workqueue.RateLimitingInterface

	// clusterLister is able to list/get Clusters from a shared informer's
	// store.
	clusterLister listersv1alpha1.BrokerClusterLister
	// clusterListerSynced returns true if the Cluster shared informer has
	// synced at least once.
	clusterListerSynced cache.InformerSynced
	// clusterUpdater implements control logic for updating Cluster
	// statuses. Implemented as an interface to enable testing.
	clusterUpdater clusterUpdaterInterface

	// statefulSetLister is able to list/get StatefulSets from a shared
	// informer's store.
	statefulSetLister appslisters.StatefulSetLister
	// statefulSetListerSynced returns true if the StatefulSet shared informer
	// has synced at least once.
	statefulSetListerSynced cache.InformerSynced
	// statefulSetControl enables control of StatefulSets associated with
	// Clusters.
	statefulSetControl StatefulSetControlInterface

	// podLister is able to list/get Pods from a shared
	// informer's store.
	podLister corelisters.PodLister
	// podListerSynced returns true if the Pod shared informer
	// has synced at least once.
	podListerSynced cache.InformerSynced
	// podControl enables control of Pods associated with
	// Clusters.
	podControl PodControlInterface

	// serviceLister is able to list/get Services from a shared informer's
	// store.
	serviceLister corelisters.ServiceLister
	// serviceListerSynced returns true if the Service shared informer
	// has synced at least once.
	serviceListerSynced cache.InformerSynced
	// serviceControl enables control of Services associated with Clusters.
	serviceControl ServiceControlInterface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewBrokerController creates a new BrokerController.
func NewBrokerController(
	opConfig operatoropts.OperatorOpts,
	opClient clientset.Interface,
	kubeClient kubernetes.Interface,
	clusterInformer informersv1alpha1.BrokerClusterInformer,
	statefulSetInformer appsinformers.StatefulSetInformer,
	podInformer coreinformers.PodInformer,
	serviceInformer coreinformers.ServiceInformer,
	resyncPeriod time.Duration,
	namespace string,
) *BrokerController {
	opscheme.AddToScheme(scheme.Scheme) // TODO: This shouldn't be done here I don't think.

	// Create event broadcaster.
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	m := BrokerController{
		opConfig: opConfig,

		opClient:   opClient,
		kubeClient: kubeClient,

		clusterLister:       clusterInformer.Lister(),
		clusterListerSynced: clusterInformer.Informer().HasSynced,
		clusterUpdater:      newClusterUpdater(opClient, clusterInformer.Lister()),

		serviceLister:       serviceInformer.Lister(),
		serviceListerSynced: serviceInformer.Informer().HasSynced,
		serviceControl:      NewRealServiceControl(kubeClient, serviceInformer.Lister()),

		statefulSetLister:       statefulSetInformer.Lister(),
		statefulSetListerSynced: statefulSetInformer.Informer().HasSynced,
		statefulSetControl:      NewRealStatefulSetControl(kubeClient, statefulSetInformer.Lister()),

		podLister:       podInformer.Lister(),
		podListerSynced: podInformer.Informer().HasSynced,
		podControl:      NewRealPodControl(kubeClient, podInformer.Lister()),

		queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "brokercluster"),
		recorder: recorder,
	}

	clusterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: m.enqueueCluster,
		UpdateFunc: func(old, new interface{}) {
			m.enqueueCluster(new)
		},
		DeleteFunc: func(obj interface{}) {
			cluster, ok := obj.(*v1alpha1.BrokerCluster)
			if ok {
				m.onClusterDeleted(cluster.Name)
			}
		},
	})

	statefulSetInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: m.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newStatefulSet := new.(*apps.StatefulSet)
			oldStatefulSet := old.(*apps.StatefulSet)
			if newStatefulSet.ResourceVersion == oldStatefulSet.ResourceVersion {
				return
			}

			// If cluster is ready ...
			if newStatefulSet.Status.ReadyReplicas == newStatefulSet.Status.Replicas {
				clusterName, ok := newStatefulSet.Labels[constants.BrokerClusterLabel]
				if ok {
					m.onClusterReady(clusterName)
				}
			}
			m.handleObject(new)
		},
		DeleteFunc: m.handleObject,
	})

	return &m
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (m *BrokerController) Run(ctx context.Context, threadiness int) {
	defer utilruntime.HandleCrash()
	defer m.queue.ShutDown()

	glog.Info("Starting Cluster controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for Cluster controller informer caches to sync")
	if !controllerutils.WaitForCacheSync("rocketmq cluster", ctx.Done(),
		m.clusterListerSynced,
		m.statefulSetListerSynced,
		m.podListerSynced,
		m.serviceListerSynced) {
		return
	}

	glog.Info("Starting Cluster controller workers")
	// Launch two workers to process Foo resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(m.runWorker, time.Second, ctx.Done())
	}

	glog.Info("Started Cluster controller workers")
	defer glog.Info("Shutting down Cluster controller workers")
	<-ctx.Done()
}

// worker runs a worker goroutine that invokes processNextWorkItem until the
// controller's queue is closed.
func (m *BrokerController) runWorker() {
	for m.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (m *BrokerController) processNextWorkItem() bool {
	obj, shutdown := m.queue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer m.queue.Done(obj)
		key, ok := obj.(string)
		if !ok {
			m.queue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in queue but got %#v", obj))
			return nil
		}
		if err := m.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		m.queue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Cluster
// resource with the current status of the resource.
func (m *BrokerController) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name.
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	nsName := types.NamespacedName{Namespace: namespace, Name: name}

	// Get the Cluster resource with this namespace/name.
	cluster, err := m.clusterLister.BrokerClusters(namespace).Get(name)
	if err != nil {
		// The Cluster resource may no longer exist, in which case we stop processing.
		if apierrors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("brokercluster '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	// Set defaults
	cluster.EnsureDefaults()

	//todo: validate cluster, update cluster labels

	// sync cluster
	groupReplica := int(cluster.Spec.GroupReplica)
	for index := 0; index < groupReplica; index++ {
		svc, err := m.serviceLister.Services(cluster.Namespace).Get(fmt.Sprintf(cluster.Name+`-svc-%d`, index))
		// If the resource doesn't exist, we'll create it
		if apierrors.IsNotFound(err) {
			glog.V(2).Infof("Creating a new Service for cluster %q", nsName)
			svc = services.NewHeadlessService(cluster, index)
			err = m.serviceControl.CreateService(svc)
		}
		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}
		// If the Service is not controlled by this Cluster resource, we should
		// log a warning to the event recorder and return.
		if !metav1.IsControlledBy(svc, cluster) {
			msg := fmt.Sprintf(MessageResourceExists, "Service", svc.Namespace, svc.Name)
			m.recorder.Event(cluster, corev1.EventTypeWarning, ErrResourceExists, msg)
			return errors.New(msg)
		}
	}

	// sync statefulsets
	membersPerGroup := int(cluster.Spec.MembersPerGroup)
	if membersPerGroup == 0 {
		utilruntime.HandleError(fmt.Errorf("invalid membersPerGroup %s", membersPerGroup))
		return nil
	}

	if cluster.Spec.AllMaster {
		cluster.Spec.MembersPerGroup = 1
	}
	/*membersPerGroup := 1
	if cluster.Spec.ClusterMode != "ALL-MASTER" && cluster.Spec.MembersPerGroup != nil{
		membersPerGroup := int(cluster.Spec.MembersPerGroup)
	}*/

	readyGroups := 0
	readyMembers := 0
	for index := 0; index < groupReplica; index++ {
		ss, err := m.statefulSetLister.StatefulSets(cluster.Namespace).Get(fmt.Sprintf(cluster.Name+`-%d`, index))
		// If the resource doesn't exist, we'll create it
		if apierrors.IsNotFound(err) {
			glog.V(2).Infof("Creating a new StatefulSet for cluster %q", nsName)
			ss = statefulsets.NewStatefulSet(cluster, index)
			err = m.statefulSetControl.CreateStatefulSet(ss)
		}
		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}
		// If the StatefulSet is not controlled by this Cluster resource, we
		// should log a warning to the event recorder and return.
		if !metav1.IsControlledBy(ss, cluster) {
			msg := fmt.Sprintf(MessageResourceExists, "StatefulSet", ss.Namespace, ss.Name)
			m.recorder.Event(cluster, corev1.EventTypeWarning, ErrResourceExists, msg)
			return fmt.Errorf(msg)
		}
		// todo: Upgrade the required component resources for the current Operator version.

		// If this number of the members on the Cluster does not equal the
		// current desired replicas on the StatefulSet, we should update the
		// StatefulSet resource.

		if cluster.Spec.MembersPerGroup != *ss.Spec.Replicas {
			glog.V(4).Infof("Updating %q: membersPerGroup=%d statefulSetReplicas=%d",
				nsName, cluster.Spec.MembersPerGroup, ss.Spec.Replicas)
			old := ss.DeepCopy()
			ss = statefulsets.NewStatefulSet(cluster, index)
			if err := m.statefulSetControl.Patch(old, ss); err != nil {
				// Requeue the item so we can attempt processing again later.
				// This could have been caused by a temporary network failure etc.
				return err
			}
		}
		if ss.Status.ReadyReplicas == ss.Status.Replicas {
			readyGroups++
			readyMembers += int(ss.Status.Replicas)
		}
	}

	// Finally, we update the status block of the Cluster resource to
	// reflect the current state of the world.
	err = m.updateClusterStatus(cluster, readyGroups, readyMembers)
	if err != nil {
		return err
	}

	m.recorder.Event(cluster, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// updateClusterStatusForSS updates Cluster statuses based on changes to their associated StatefulSets.
func (m *BrokerController) updateClusterStatus(cluster *v1alpha1.BrokerCluster, readyGroups int, readyMembers int) error {
	glog.V(4).Infof("%s/%s: readyGroups=%d, readyMembers=%d",
		cluster.Namespace, cluster.Name, readyGroups, readyMembers)

	status := cluster.Status.DeepCopy()
	_, condition := clusterutil.GetClusterCondition(&cluster.Status, v1alpha1.BrokerClusterReady)
	if condition == nil {
		condition = &v1alpha1.BrokerClusterCondition{Type: v1alpha1.BrokerClusterReady}
	}
	if readyGroups == int(cluster.Spec.GroupReplica) && readyMembers == int(cluster.Spec.GroupReplica)*int(cluster.Spec.MembersPerGroup) {
		condition.Status = corev1.ConditionTrue
	} else {
		condition.Status = corev1.ConditionFalse
	}

	if updated := clusterutil.UpdateClusterCondition(status, condition); updated {
		return m.clusterUpdater.UpdateClusterStatus(cluster.DeepCopy(), status)
	}
	return nil
}

// enqueueCluster takes a Cluster resource and converts it into a
// namespace/name string which is then put onto the work queue. This method
// should *not* be passed resources of any type other than Cluster.
func (m *BrokerController) enqueueCluster(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	m.queue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Foo resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (m *BrokerController) handleObject(obj interface{}) {
	object, ok := obj.(metav1.Object)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}

	glog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Cluster, we should not do
		// anything more with it.
		if ownerRef.Kind != v1alpha1.ClusterCRDResourceKind {
			return
		}

		cluster, err := m.clusterLister.BrokerClusters(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of Cluster '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		m.enqueueCluster(cluster)
		return
	}
}

func (m *BrokerController) onClusterReady(clusterName string) {
	glog.V(2).Infof("Cluster %s ready", clusterName)
}

func (m *BrokerController) onClusterDeleted(clusterName string) {
	glog.V(2).Infof("Cluster %s deleted", clusterName)
}
