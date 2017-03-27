package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sam701/grpc-kv/kv"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:      "set",
			Usage:     "Set value",
			ArgsUsage: "KEY VALUE",
			Action:    setKeyValue,
		},
		{
			Name:      "get",
			Usage:     "Get value",
			ArgsUsage: "KEY",
			Action:    getValue,
		},
	}
	app.Before = func(ctx *cli.Context) error {
		client = createClient()
		return nil
	}
	app.Run(os.Args)

}

var client kv.KeyValueStoreClient

func createClient() kv.KeyValueStoreClient {
	conn, err := grpc.Dial("localhost:12000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return kv.NewKeyValueStoreClient(conn)
}

func setKeyValue(ctx *cli.Context) error {
	key := ctx.Args().First()
	value := ctx.Args().Get(1)
	_, err := client.Set(context.Background(), &kv.KeyValue{key, value})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return nil
}
func getValue(ctx *cli.Context) error {
	key := ctx.Args().First()
	kv, err := client.Get(context.Background(), &kv.Key{key})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	fmt.Println(kv.Value)
	return nil
}
