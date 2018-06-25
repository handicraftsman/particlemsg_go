package main

import (
	"crypto/tls"
	"fmt"
	"log"

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

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	go srv.Start(host, port, true, &tls.Config{Certificates: []tls.Certificate{cer}}, func(mi particlemsg.MessageInfo) {
		fmt.Printf("%s: %v\n", mi.From, mi.Msg)
	})

	particlemsg.LoadPlugins(host, port, clients)

	<-srv.DoneChan
}
