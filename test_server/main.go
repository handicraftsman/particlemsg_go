package main

import (
	"fmt"
	"time"

	"github.com/handicraftsman/particlemsg_go"
)

func main() {
	srv := particlemsg.NewServer()

	host, port, crt, key := particlemsg.GetServerConfig()

	clients := particlemsg.LoadClientConfig("./clients.json")
	srv.LoadClientConfig(clients)

	srv.Blocked = true

	go srv.Start(
		host,
		port,
		true,
		particlemsg.GetBasicSSLConfig(particlemsg.GetSSLCertFromFiles(crt, key)),
		func(mi particlemsg.MessageInfo) {
			fmt.Printf("%s: %v\n", mi.From, mi.Msg)
			if mi.Msg.Type == "_registered" {
				who := mi.Msg.Data["Who"].(string)
				if who == "core" {
					srv.Blocked = false
					particlemsg.LoadPlugins("127.0.0.1", port, crt, key, clients, true)
				}
				time.Sleep(time.Second)
				srv.Broadcast(particlemsg.Message{
					Type: "newPlugin",
					Data: map[string]interface{}{
						"Who": who,
					},
				})
			} else if mi.Msg.Type == "_subscribe" {
				time.Sleep(time.Second)
				srv.BroadcastToSubscribed(particlemsg.Message{
					Type: mi.Msg.Data["Type"].(string),
					Data: mi.Msg.Data["Pattern"].(map[string]interface{}),
				})
			}
		})

	core, err := particlemsg.FindClientInfo(clients, "core")
	if err != nil {
		panic(err)
	}
	particlemsg.LoadPlugin("127.0.0.1", port, crt, key, core)

	<-srv.DoneChan
}
