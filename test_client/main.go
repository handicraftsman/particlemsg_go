package main

import (
	"fmt"
	"os"

	particlemsg "github.com/handicraftsman/particlemsg_go"
)

func main() {
	c := particlemsg.NewClient(os.Getenv("PMSG_NAME"))

	fmt.Println(os.Getenv("ASDF"))

	go c.ConnectFromEnv(
		particlemsg.GetBasicSSLConfig(particlemsg.GetSSLCertFromFiles(os.Getenv("PMSG_SSL_CERT"), os.Getenv("PMSG_SSL_KEY"))),
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
