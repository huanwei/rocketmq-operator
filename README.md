# RocketMQ Operator

The RoecketMQ Operator creates, operates and scales self-healing RocketMQ clusters on Kubernetes.

** This is currently in v1alpha. **

** By default, RoecketMQ Operator creates the storage with EmptyDir.

** Use this operator in production with Kubernetes [storageClass](docs/storageClass.md) **

## Getting started

See the [tutorial](docs/tutorial.md) for a quick-start guide.

## Features

The RoecketMQ Operator provides following key features:

- Create, operate and scale self-healing RocketMQ clusters on Kubernetes.
- Other featues are WIP.

## Prerequisites 

* Kubernetes 1.9+, since StatefulSets are stable (GA) in 1.9.

## Contributing 

`rocketmq-operator` is an open source project. Welcome to submit PR.

## License

Copyright The Kubernetes Authors.

`rocketmq-operator` is licensed under the Apache License 2.0. 

See [LICENSE](LICENSE) for more details.