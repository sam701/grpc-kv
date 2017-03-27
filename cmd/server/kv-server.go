package main

import (
	"log"
	"net"

	"github.com/sam701/grpc-kv/kv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:12000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	kv.RegisterKeyValueStoreServer(s, &server{make(map[string]string)})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("done")
}

type server struct {
	data map[string]string
}

func (s *server) Set(ctx context.Context, in *kv.KeyValue) (*kv.Empty, error) {
	s.data[in.Key] = in.Value
	return &kv.Empty{}, nil
}

func (s *server) Get(ctx context.Context, in *kv.Key) (*kv.KeyValue, error) {
	return &kv.KeyValue{in.Key, s.data[in.Key]}, nil
}
