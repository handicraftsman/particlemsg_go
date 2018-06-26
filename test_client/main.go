package main

import (
	"fmt"
	"os"

	particlemsg "github.com/handicraftsman/particlemsg_go"
)

func main() {
	c := particlemsg.NewClient(os.Getenv("PMSG_NAME"))

	go c.ConnectFromEnv(
		particlemsg.GetBasicSSLConfig(particlemsg.GetSSLCertFromFiles("./client.crt", "./client.key")),
		func(mi particlemsg.MessageInfo) {
			fmt.Printf("%s: %v\n", mi.From, mi.Msg)
			if mi.Msg.Type == "_registered" {
				mi.CConn.SendTo(os.Getenv("PMSG_NAME"), particlemsg.Message{
					Type: "foo",
					Data: map[string]interface{}{
						"foo": "bar",
						"baz": "quux",
					},
				})
			}
		})

	<-c.DoneChan
}
