package main

import (
	"fmt"
	"os"

	"github.com/handicraftsman/particlemsg_go"
)

var (
	host = os.Getenv("PMSG_HOST")
	port = os.Getenv("PMSG_PORT")
	crt  = os.Getenv("PMSG_SSL_CERT")
	key  = os.Getenv("PMSG_SSL_KEY")
)

func main() {
	srv := particlemsg.NewServer()

	clients := particlemsg.LoadClientConfig("./clients.json")
	srv.LoadClientConfig(clients)

	go srv.Start(
		host,
		port,
		true,
		particlemsg.GetBasicSSLConfig(particlemsg.GetSSLCertFromFiles(crt, key)),
		func(mi particlemsg.MessageInfo) {
			fmt.Printf("%s: %v\n", mi.From, mi.Msg)
		})

	particlemsg.LoadPlugins("127.0.0.1", port, crt, key, clients)

	<-srv.DoneChan
}
