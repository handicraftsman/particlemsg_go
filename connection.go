package particlemsg

import (
	"crypto/tls"
	"net"
	"sync"
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
	Conn         *net.Conn
	Name         string
	SubscribedTo sync.Map
}

// SubPair - identifies SubInfos in ServerConnection.SubscribedTo
type SubPair struct {
	Type string
	ID   string
}

// SubInfo - describes subscriptions
type SubInfo struct {
	ID      string
	Type    string
	Pattern map[string]interface{}
	Once    bool
}

// NewServerConnection - creates a new ServerConnection
func NewServerConnection(conn *net.Conn) *ServerConnection {
	return &ServerConnection{
		Conn:         conn,
		Name:         "_unregistered",
		SubscribedTo: sync.Map{},
	}
}

// SendMessage - sends a message
func (c *ServerConnection) SendMessage(msg Message) {
	jl := MessageToJSON(msg)
	(*c.Conn).Write([]byte(jl + "\r\n"))
}

// Subscribe - subscribes to the given pattern
func (c *ServerConnection) Subscribe(t string, d *map[string]interface{}, id string, once bool) {
	c.SubscribedTo.Store(SubPair{Type: t, ID: id}, SubInfo{ID: id, Type: t, Pattern: *d, Once: once})
}

// Unsubscribe - unsubscribes from the given pattern
func (c *ServerConnection) Unsubscribe(id string) bool {
	var t string
	ok := false
	c.SubscribedTo.Range(func(_k, _ interface{}) bool {
		k := _k.(SubPair)
		if k.ID == id {
			t = k.Type
			return false
		}
		return true
	})
	if ok {
		c.SubscribedTo.Delete(SubPair{Type: t, ID: id})
		return true
	}
	return false
}

// IsSubscribedTo - checks if connection is subscribed to the given message
func (c *ServerConnection) IsSubscribedTo(what Message) bool {
	ok := false
	del := false
	delpair := SubPair{}
	c.SubscribedTo.Range(func(_, _v interface{}) bool {
		v := _v.(SubInfo)
		if what.Type == v.Type && isSupersetOf(what.Data, v.Pattern) {
			if v.Once {
				delpair = SubPair{ID: v.ID, Type: v.Type}
			}
			ok = true
			return false
		}
		return true
	})
	if del {
		c.SubscribedTo.Delete(delpair)
	}
	return ok
}
