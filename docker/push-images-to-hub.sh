#!/bin/bash

for i in client server; do
  docker push sam701/kv-$i
done