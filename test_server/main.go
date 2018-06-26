package main

import (
	"fmt"

	"github.com/handicraftsman/particlemsg_go"
)

var (
	host         = "0.0.0.0"
	port         = "5050"
	allowedHosts = []string{"127.0.0.1"}
)

func main() {
	srv := particlemsg.NewServer()

	clients := particlemsg.LoadClientConfig("./clients.json")
	srv.LoadClientConfig(clients)

	go srv.Start(
		host,
		port,
		true,
		particlemsg.GetBasicSSLConfig(particlemsg.GetSSLCertFromFiles("./server.crt", "./server.key")),
		func(mi particlemsg.MessageInfo) {
			fmt.Printf("%s: %v\n", mi.From, mi.Msg)
		})

	particlemsg.LoadPlugins(host, port, clients)

	<-srv.DoneChan
}
