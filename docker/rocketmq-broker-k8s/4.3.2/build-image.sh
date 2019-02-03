#!/bin/bash

docker build -t huanwei/rocketmq-broker-k8s:4.3.2 .
docker push huanwei/rocketmq-broker-k8s:4.3.2