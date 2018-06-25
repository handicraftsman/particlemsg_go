package particlemsg

import (
	"encoding/json"
)

// Message - struct representing a message
type Message struct {
	Type string
	Data map[string]interface{}
}

// MessageInfo - stores message pointer, connection pointer and sender name (_server if message is sent by the server itself)
type MessageInfo struct {
	Msg   *Message
	CConn *ClientConnection
	SConn *ServerConnection
	From  string
}

// JSONToMessage - converts json data into a message
func JSONToMessage(jl string) Message {
	var msg Message
	err := json.Unmarshal([]byte(jl), &msg)
	if err != nil {
		panic(err)
	}
	return msg
}

// MessageToJSON - converts message into json
func MessageToJSON(msg Message) string {
	var jl string
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	jl = string(b)
	return jl
}
