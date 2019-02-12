#!/bin/bash

docker build -t rocketmqinc/rocketmq-operator:0.1.2-dev .
docker push rocketmqinc/rocketmq-operator:0.1.2-dev