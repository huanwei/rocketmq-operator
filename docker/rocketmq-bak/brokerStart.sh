#!/bin/bash

./brokerGenConfig.sh
#./mqbroker -n rocketmq-namesrv-service.rocketmq-operator:9876
#./mqbroker -n 192.168.196.49:9876;192.168.122.177:9876
./mqbroker -n $NAMESRV_ADDRESS