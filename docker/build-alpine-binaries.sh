#!/bin/bash

cd ..

docker run --name grpc-kv --rm \
  -v "$PWD":/go/src/grpc-kv -w /go/src/grpc-kv \
  sam701/golang-and-dependencies \
  sh -c 'go get -v ./...; go build ./cmd/kv-client; go build ./cmd/kv-server'
  
mv kv-server docker/server/
mv kv-client docker/client/
