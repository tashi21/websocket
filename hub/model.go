package hub

import "github.com/gorilla/websocket"

// maintain the set of active connections and topics
type Hub struct {
	Topics    map[string]map[*Cli]bool // registered topics
	Broadcast chan Msg                 // channel of inbound messages from the subscribers
	Reg       chan *Cli                // channel to register client to given topics
	UnReg     chan *Cli                // channel to unregister client from given topics when they close connection
	clients   int64                    // number of clients connected to the hub
}

// client struct
type Cli struct {
	Conn       *websocket.Conn // websocket connection
	Registered bool            // is client registered to hub
	Send       chan Msg        // channel of messages to send to the client from the hub
	Topic      string          // topic subscribed to by the client
	tid        string          // trace id
}

// message type to send and recieve from the client
type Msg struct {
	Topic string                 `json:"topic"`   // topic to which the message is sent
	Msg   map[string]interface{} `json:"message"` // message to be sent
}
