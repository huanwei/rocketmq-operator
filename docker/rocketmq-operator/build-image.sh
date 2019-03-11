#!/bin/bash

#docker build -t huanwei/rocketmq-operator:0.1.2 .
#docker push huanwei/rocketmq-operator:0.1.2

docker build -t rocketmqinc/rocketmq-operator:0.1.2-dev .
docker push rocketmqinc/rocketmq-operator:0.1.2-dev