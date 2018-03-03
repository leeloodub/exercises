package main

import (
	"flag"
	"log"
	"net"
)

var (
	eventTargetAddr   string
	clientManagerAddr string
	endpoint          *Endpoint
	clientManager     *ClientManager
)

const (
	protocol = "tcp"
)

func init() {
	flag.StringVar(&eventTargetAddr, "event_target_addr", ":9090", "address/port to listen on")
	flag.StringVar(&clientManagerAddr, "client_manager_addr", ":9099", "address/port to listen on")
}

func main() {
	etl, err := net.Listen("tcp", eventTargetAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer etl.Close()
	endpoint = NewEndpoint(etl, startSeq)
	go endpoint.ListenEvents()

	cml, err := net.Listen("tcp", clientManagerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer cml.Close()
	clientManager = NewClientManager()
	endpoint.processor = clientManager.processMessage

	clientManager.ManageClients(cml)
}
