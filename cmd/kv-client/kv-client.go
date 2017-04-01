package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sam701/grpc-kv/kv"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kv-service",
			Usage:  "Consul ervice name",
			EnvVar: "KV_SERVICE",
		},
		cli.StringFlag{
			Name:   "kv-server",
			Usage:  "KV-Server `HOST:PORT` to connect to (overwrites service)",
			EnvVar: "KV_SERVER",
		},
		cli.StringFlag{
			Name:   "listen",
			Usage:  "Web server `HOST:PORT` to listen on",
			Value:  ":80",
			EnvVar: "LISTEN",
		},
		cli.StringFlag{
			Name:   "id",
			Usage:  "Server `ID`",
			EnvVar: "ID",
		},
	}
	app.Action = runWebServer
	app.Run(os.Args)

}

func runWebServer(ctx *cli.Context) error {
	serviceName = ctx.String("kv-service")
	serverAddr = ctx.String("kv-server")

	client = createClient()

	r := mux.NewRouter()

	id := ctx.String("id")
	if id == "" {
		id = "client-" + strconv.Itoa(int(rand.Uint32()))
	}
	r.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(id))
	})

	r.Methods("GET").Path("/{key}").HandlerFunc(getHandler)
	r.Methods("PUT").Path("/{key}/{value}").HandlerFunc(setHandler)

	http.ListenAndServe(ctx.String("listen"), r)
	return nil
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	kv, err := client.Get(context.Background(), &kv.Key{key})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if kv.Value == "" {
		http.Error(w, "Not found", 404)
	} else {
		w.Write([]byte(kv.Value))
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	_, err := client.Set(context.Background(), &kv.KeyValue{key, value})
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.WriteHeader(204)
	}
}

var client kv.KeyValueStoreClient

func createClient() kv.KeyValueStoreClient {
	conn, err := grpc.Dial(getAddr(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return kv.NewKeyValueStoreClient(conn)
}
