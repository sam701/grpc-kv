package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc/naming"

	consul "github.com/hashicorp/consul/api"
)

type addrResolver struct{}

func (s *addrResolver) Resolve(target string) (naming.Watcher, error) {
	log.Println("Resolving:", target)
	return newAddrWatcher(target), nil
}

type addrWatcher struct {
	target string

	currentAddr string
	updatesChan chan []*naming.Update
}

func newAddrWatcher(target string) *addrWatcher {
	r := &addrWatcher{
		target:      target,
		updatesChan: make(chan []*naming.Update),
	}
	go r.watchAddressChange()
	return r
}

func (s *addrWatcher) watchAddressChange() {
	cfg := consul.DefaultConfig()
	cfg.Address = consulAddr

	c, err := consul.NewClient(cfg)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	var waitIndex uint64
	for {
		services, meta, err := c.Catalog().Service(s.target, "", &consul.QueryOptions{
			WaitIndex: waitIndex,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		waitIndex = meta.LastIndex

		if len(services) > 0 {
			se := services[0]
			s.setCurrentAddress(fmt.Sprintf("%s:%d", se.ServiceAddress, se.ServicePort))
		}
	}
}

func (s *addrWatcher) Next() ([]*naming.Update, error) {
	ups := <-s.updatesChan
	return ups, nil
}

func (s *addrWatcher) Close() {}

func (s *addrWatcher) setCurrentAddress(addr string) {
	if addr == s.currentAddr {
		log.Println("No changes for current address:", addr)
		return
	}

	updates := make([]*naming.Update, 0, 2)

	updates = append(updates, &naming.Update{
		Op:   naming.Add,
		Addr: addr,
	})

	if s.currentAddr != "" {
		updates = append(updates, &naming.Update{
			Op:   naming.Delete,
			Addr: s.currentAddr,
		})
	}

	s.currentAddr = addr

	log.Println("Setting current address:", s.currentAddr)

	s.updatesChan <- updates
}
