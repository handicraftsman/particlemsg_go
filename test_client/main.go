package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	particlemsg "github.com/handicraftsman/particlemsg_go"
)

func main() {
	c := particlemsg.NewClient(os.Getenv("PMSG_NAME"))

	cer, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Println(err)
		return
	}

	go c.ConnectFromEnv(&tls.Config{Certificates: []tls.Certificate{cer}, InsecureSkipVerify: true}, func(mi particlemsg.MessageInfo) {
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
