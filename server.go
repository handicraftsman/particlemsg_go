package particlemsg

import (
	"bufio"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// Server - struct representing a server. Use it when you want to listen for plugin connections
type Server struct {
	DoneChan    chan bool
	Connections map[string]*ServerConnection
	MsgChan     chan MessageInfo
	Keys        map[string]string
}

// NewServer - creates a new server
func NewServer() *Server {
	return &Server{
		DoneChan:    make(chan bool),
		Connections: make(map[string]*ServerConnection),
		MsgChan:     make(chan MessageInfo, 16),
		Keys:        make(map[string]string),
	}
}

// LoadClientConfig - loads ClientConfig into the server
func (s *Server) LoadClientConfig(clients *ClientConfig) {
	for _, client := range *clients {
		s.Keys[client.Name] = fmt.Sprintf("%x", sha256.Sum256([]byte(client.Key)))
	}
}

// Start - starts the server
func (s *Server) Start(host, port string, requireKeys bool, tcfg *tls.Config, f func(MessageInfo)) {
	l, err := tls.Listen("tcp", host+":"+port, tcfg)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	fmt.Println("Listening on " + host + ":" + port)
	for {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Panic: %v\n", err)
			}
		}()

		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go s.handleClient(conn, requireKeys, f)
	}
}

func (s *Server) handleClient(_conn net.Conn, requireKey bool, f func(MessageInfo)) {
	registered := false
	var name string
	conn := NewServerConnection(&_conn)
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

	running := true
	ponged := true
	go func() {
		upd := time.Now()
		for running {
			time.Sleep(time.Second)
			c := time.Now()
			if c.Sub(upd) >= time.Minute {
				upd = c
				if ponged {
					ponged = false
					conn.SendMessage(Message{
						Type: "_ping",
						Data: map[string]interface{}{},
					})
				} else {
					conn.SendMessage(Message{
						Type: "_quit",
						Data: map[string]interface{}{
							"Reason": "Ping Timeout",
						},
					})
					(*conn.Conn).Close()
					running = false
				}
			}
		}
	}()

	scanner := bufio.NewScanner(*conn.Conn)
	for scanner.Scan() {
		line := scanner.Text()
		msg := JSONToMessage(line)
		if msg.Type == "_pong" {
			ponged = true
		} else if msg.Type == "_register" {
			var key string
			if msg.Data["Key"] != nil {
				key = msg.Data["Key"].(string)
			}
			if !registered {
				name = msg.Data["Name"].(string)
			}
			if _, ok := s.Connections[name]; registered || ok {
				conn.SendMessage(Message{
					Type: "_alreadyRegistered",
					Data: map[string]interface{}{
						"Name": name,
					},
				})
				continue
			}
			if k, ok := s.Keys[name]; (requireKey && !ok) || (requireKey && (key != k)) {
				conn.SendMessage(Message{
					Type: "_invalidKey",
					Data: map[string]interface{}{
						"Key": key,
					},
				})
				continue
			}
			s.Connections[name] = conn
			conn.Name = name
			registered = true
			conn.SendMessage(Message{
				Type: "_registered",
				Data: map[string]interface{}{
					"Name": name,
				},
			})
		} else if registered && msg.Type == "_message" {
			to := msg.Data["To"].(string)
			message := msg.Data["Message"]
			if _, ok := s.Connections[to]; ok {
				s.Connections[to].SendMessage(Message{
					Type: "_message",
					Data: map[string]interface{}{
						"From":    name,
						"Message": message,
					},
				})
			} else {
				conn.SendMessage(Message{
					Type: "_messageError",
					Data: map[string]interface{}{
						"Reason": to + " is offline",
					},
				})
			}
		} else if msg.Type == "_quit" {
			delete(s.Connections, name)
			conn.SendMessage(Message{
				Type: "_quit",
				Data: map[string]interface{}{
					"Reason": "Client Quit",
				},
			})
			(*conn.Conn).Close()
			running = false
		}
		if registered || msg.Type == "_quit" || msg.Type == "_register" || msg.Type == "_pong" {
			go f(MessageInfo{
				Msg:   &msg,
				SConn: conn,
				From:  name,
			})
		}
	}
	delete(s.Connections, name)
	(*conn.Conn).Close()
	running = false
	s.DoneChan <- true
}
