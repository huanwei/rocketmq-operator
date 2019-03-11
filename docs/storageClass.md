# Tutorial on Kubernetes StorageClass

For production scenarios, we often need to dynamically provision persistent volumes and persistent volume claims, so we need to use Kubernetes StorageClass.

For example, we use StorageClass like below:

```
[root@k8s-master ~]# kubectl get storageclass
NAME                     PROVISIONER          AGE
demo-nfs-storage         harmonycoud.cn/nfs   56d
mongo-demo-nfs-storage   harmonycoud.cn/nfs   67d

[root@k8s-master ~]# kubectl get storageclass demo-nfs-storage  -o yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  creationTimestamp: 2019-01-14T12:35:02Z
  name: demo-nfs-storage
  resourceVersion: "30268939"
  selfLink: /apis/storage.k8s.io/v1/storageclasses/demo-nfs-storage
  uid: d0125b5b-17f8-11e9-8963-005056bc2383
provisioner: harmonycoud.cn/nfs
reclaimPolicy: Delete

```

So with Kubernetes StorageClass our RocketMQ Operator only needs to know the `storageClassName`, and all the RocketMQ broker clusters PV and PVC would be provisioned automatically, like following:
```
[root@k8s-master ~]# kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                                               STORAGECLASS       REASON    AGE
pvc-920d2bba-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerlogs-mybrokercluster-0-0    demo-nfs-storage             3h
pvc-920ead52-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerstore-mybrokercluster-0-0   demo-nfs-storage             3h
pvc-9225c2cf-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerlogs-mybrokercluster-1-0    demo-nfs-storage             3h
pvc-922c1b7f-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerstore-mybrokercluster-1-0   demo-nfs-storage             3h
pvc-9631143d-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerstore-mybrokercluster-1-1   demo-nfs-storage             3h
pvc-96331246-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerlogs-mybrokercluster-1-1    demo-nfs-storage             3h
pvc-96fa5022-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerlogs-mybrokercluster-0-1    demo-nfs-storage             3h
pvc-96fd203d-43da-11e9-9df9-005056bc2383   5Gi        RWO            Delete           Bound     rocketmq-operator/brokerstore-mybrokercluster-0-1   demo-nfs-storage             3h


[root@k8s-master ~]# kubectl get pvc -n rocketmq-operator
NAME                              STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS       AGE
brokerlogs-mybrokercluster-0-0    Bound     pvc-920d2bba-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h
brokerlogs-mybrokercluster-0-1    Bound     pvc-96fa5022-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h
brokerlogs-mybrokercluster-1-0    Bound     pvc-9225c2cf-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h
brokerlogs-mybrokercluster-1-1    Bound     pvc-96331246-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h
brokerstore-mybrokercluster-0-0   Bound     pvc-920ead52-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h
brokerstore-mybrokercluster-0-1   Bound     pvc-96fd203d-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h
brokerstore-mybrokercluster-1-0   Bound     pvc-922c1b7f-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h
brokerstore-mybrokercluster-1-1   Bound     pvc-9631143d-43da-11e9-9df9-005056bc2383   5Gi        RWO            demo-nfs-storage   3h

```

