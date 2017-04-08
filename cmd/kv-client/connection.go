package main

import (
	"log"

	"google.golang.org/grpc"

	"github.com/sam701/grpc-kv/kv"
)

var (
	client kv.KeyValueStoreClient
)

func createClient(service, server string) kv.KeyValueStoreClient {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	if service != "" {
		opts = append(opts, grpc.WithBalancer(grpc.RoundRobin(&addrResolver{})))
		server = service
	}
	conn, err := grpc.Dial(server, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return kv.NewKeyValueStoreClient(conn)
}
