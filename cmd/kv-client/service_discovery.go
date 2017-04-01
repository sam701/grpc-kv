package main

import (
	"log"

	"fmt"

	consul "github.com/hashicorp/consul/api"
)

var (
	serviceName string
	serverAddr  string
)

func getAddr() string {
	if serverAddr != "" {
		return serverAddr
	}

	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	services, _, err := c.Catalog().Service(serviceName, "", &consul.QueryOptions{})
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	s := services[0]
	a := fmt.Sprintf("%s:%d", s.ServiceAddress, s.ServicePort)
	log.Println("addr:", a)
	return a
}
