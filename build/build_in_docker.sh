#!/bin/bash
if [ -z "$GOPATH" ]; then
  GOPATH="$HOME/Downloads/gopath"
fi
SRC_HOME_RELATIVE=src/github.com/huanwei/rocketmq-operator
TARGET_FILE_RELATIVE=${SRC_HOME_RELATIVE}/docker/rocketmq-operator/rocketmq-operator
TARGET_FILE=${GOPATH}/${TARGET_FILE_RELATIVE}
TARGET_FILE_IN_DOCKER=/gopath/${TARGET_FILE_RELATIVE}

echo "Removing old bin file..."
rm -rf ${TARGET_FILE}

echo "GOPATH is $GOPATH"
echo  "Building rocketmq-operator..."
docker run --rm -v "$GOPATH":/gopath -w /gopath/${SRC_HOME_RELATIVE}/ -e GOPATH=/gopath golang:1.10 sh -c 'cd cmd/rocketmq-operator && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../docker/rocketmq-operator/rocketmq-operator'
chmod +x ${TARGET_FILE}
echo  "Building rocketmq-operator complete"
stat ${TARGET_FILE}