package particlemsg

import (
	"bufio"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Client - a struct representing a client. Use it to connect to a server
type Client struct {
	DoneChan chan bool
	Name     string
	Conn     *ClientConnection
	Running  bool
}

// ClientInfo - a struct represening client info loaded from json config
type ClientInfo struct {
	Name      string
	Path      string
	Key       string
	UnsafeSSL bool
	DoNotLoad bool
	Env       []string
}

// ClientConfig - an array of ClientInfos
type ClientConfig []ClientInfo

// LoadClientConfig - loads ClientConfig from a json file
func LoadClientConfig(path string) *ClientConfig {
	var clients ClientConfig
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &clients)
	if err != nil {
		panic(err)
	}
	return &clients
}

// NewClient - creates a new client
func NewClient(name string) *Client {
	return &Client{
		DoneChan: make(chan bool),
		Name:     name,
		Conn:     nil,
		Running:  false,
	}
}

// ConnectFromEnv - connects to the server using PMSG_HOST, PMSG_PORT and PMSG_KEY environment variables
func (c *Client) ConnectFromEnv(tcfg *tls.Config, f func(MessageInfo)) {
	c.Connect(os.Getenv("PMSG_HOST"), os.Getenv("PMSG_PORT"), os.Getenv("PMSG_KEY"), tcfg, f)
}

// Connect - connects to the server. Callback is used to notify you about new messages.
func (c *Client) Connect(host, port, key string, tcfg *tls.Config, f func(MessageInfo)) {
	c.Running = true

	_conn, err := tls.Dial("tcp", host+":"+port, tcfg)
	if err != nil {
		panic(err)
	}

	conn := NewClientConnection(_conn)
	conn.Name = c.Name
	c.Conn = conn
	defer func() {
		if err := recover(); err != nil {
			conn.SendMessage(Message{
				Type: "_error",
				Data: map[string]interface{}{
					"Reason": "Panic: " + fmt.Sprintf("%v", err),
				},
			})
		}
	}()

	go func() {
		conn.SendMessage(Message{
			Type: "_register",
			Data: map[string]interface{}{
				"Name": c.Name,
				"Key":  fmt.Sprintf("%x", sha256.Sum256([]byte(key))),
			},
		})
	}()

	scanner := bufio.NewScanner(_conn)
	for scanner.Scan() {
		line := scanner.Text()
		msg := JSONToMessage(line)

		if msg.Type == "_ping" {
			conn.SendMessage(Message{
				Type: "_pong",
				Data: map[string]interface{}{},
			})
		}

		go f(MessageInfo{
			Msg:   &msg,
			CConn: conn,
			From:  "_server",
		})

		if msg.Type == "_message" {
			w := msg.Data["From"].(string)
			m := msg.Data["Message"].(map[string]interface{})
			t := m["Type"].(string)
			d := m["Data"].(map[string]interface{})
			go f(MessageInfo{
				Msg: &Message{
					Type: t,
					Data: d,
				},
				From: w,
			})
		}
	}

	(*_conn).Close()
	c.Running = false
	f(MessageInfo{
		Msg: &Message{
			Type: "_disconnect",
			Data: map[string]interface{}{},
		},
		CConn: conn,
		From:  "_server",
	})
	c.DoneChan <- true
}

// Subscribe - subscribes client to the given pattern
func (c *Client) Subscribe(id, t string, pattern map[string]interface{}, once bool) {
	m := Message{
		Type: "_subscribe",
		Data: map[string]interface{}{
			"ID":      id,
			"Type":    t,
			"Pattern": pattern,
			"Once":    once,
		},
	}
	c.Conn.SendMessage(m)
}

// Unsubscribe - unsubscribes client from the given pattern
func (c *Client) Unsubscribe(id string) {
	c.Conn.SendMessage(Message{
		Type: "_unsubscribe",
		Data: map[string]interface{}{
			"ID": id,
		},
	})
}
