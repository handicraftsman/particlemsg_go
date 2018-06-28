package main

import (
	"fmt"
	"os"

	particlemsg "github.com/handicraftsman/particlemsg_go"
)

func main() {
	name, crt, key := particlemsg.GetClientConfig()

	c := particlemsg.NewClient(name)

	fmt.Println(os.Getenv("ASDF"))

	go c.ConnectFromEnv(
		particlemsg.GetBasicSSLConfig(particlemsg.GetSSLCertFromFiles(crt, key)),
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
