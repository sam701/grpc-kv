#!/bin/bash

protoc kv/service.proto --go_out=plugins=grpc:.
