package particlemsg

import (
	"crypto/tls"
	"net"
)

// ClientConnection - a struct for client-side connection representation
type ClientConnection struct {
	Conn *tls.Conn
	Name string
}

// NewClientConnection - creates a new ClientConnection
func NewClientConnection(conn *tls.Conn) *ClientConnection {
	return &ClientConnection{
		Conn: conn,
		Name: "_unregistered",
	}
}

// SendMessage - sends a message
func (c *ClientConnection) SendMessage(msg Message) {
	jl := MessageToJSON(msg)
	(*c.Conn).Write([]byte(jl + "\r\n"))
}

// SendTo - sends a message to another client
func (c *ClientConnection) SendTo(to string, msg Message) {
	c.SendMessage(Message{
		Type: "_message",
		Data: map[string]interface{}{
			"To":      to,
			"Message": msg,
		},
	})
}

// ServerConnection - a struct for server-side connection representation
type ServerConnection struct {
	Conn *net.Conn
	Name string
}

// NewServerConnection - creates a new ServerConnection
func NewServerConnection(conn *net.Conn) *ServerConnection {
	return &ServerConnection{
		Conn: conn,
		Name: "_unregistered",
	}
}

// SendMessage - sends a message
func (c *ServerConnection) SendMessage(msg Message) {
	jl := MessageToJSON(msg)
	(*c.Conn).Write([]byte(jl + "\r\n"))
}
