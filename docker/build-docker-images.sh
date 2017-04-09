#!/bin/bash

(
  cd ..

  docker run --name grpc-kv --rm \
    -v "$PWD":/go/src/github.com/sam701/grpc-kv -w /go/src/github.com/sam701/grpc-kv \
    sam701/golang-and-dependencies \
    sh -c 'go build ./cmd/kv-client; go build ./cmd/kv-server'
    
  mv kv-server docker/server/
  mv kv-client docker/client/
)

for i in client server; do
  (
    cd $i
    docker build -t sam701/kv-$i .
  )
done