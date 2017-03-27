package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/sam701/grpc-kv/kv"

	"strings"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host-port",
			Usage:  "`HOST:PORT` to listen on",
			Value:  ":80",
			EnvVar: "HOST_PORT",
		},
		cli.StringFlag{
			Name:   "storage",
			Usage:  "`PATH` for key value storage",
			Value:  "/var/lib/kv-server/storage",
			EnvVar: "STORAGE",
		},
	}
	app.Action = startServer
	app.Run(os.Args)
}

func startServer(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", ctx.String("host-port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	kv.RegisterKeyValueStoreServer(s, newServer(ctx.String("storage")))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("done")
	return nil
}

type server struct {
	storagePath string
	data        map[string]string
}

func readKv(path string) map[string]string {
	data := map[string]string{}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data
		} else {
			log.Fatalln("ERROR: cannot open file:", err)
		}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, "=")
		data[parts[0]] = parts[1]
	}

	return data
}

func newServer(storagePath string) *server {
	err := os.MkdirAll(path.Dir(storagePath), 0700)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return &server{
		storagePath: storagePath,
		data:        readKv(storagePath),
	}
}

func (s *server) Set(ctx context.Context, in *kv.KeyValue) (*kv.Empty, error) {
	s.data[in.Key] = in.Value

	f, err := os.OpenFile(s.storagePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("ERROR: cannot open file:", err)
	}
	defer f.Close()
	fmt.Fprintf(f, "%s=%s\n", in.Key, in.Value)

	return &kv.Empty{}, nil
}

func (s *server) Get(ctx context.Context, in *kv.Key) (*kv.KeyValue, error) {
	return &kv.KeyValue{in.Key, s.data[in.Key]}, nil
}
