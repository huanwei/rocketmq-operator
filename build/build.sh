#!/bin/bash

rm -rf ../docker/rocketmq-operator/rocketmq-operator

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -i -o ../docker/rocketmq-operator/rocketmq-operator  ../cmd/rocketmq-operator/main.go
