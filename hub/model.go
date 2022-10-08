package hub

import "github.com/gorilla/websocket"

// maintain the set of active connections and topics
type Hub struct {
	Topics    map[string]map[*Cli]bool // registered topics
	Broadcast chan Payload             // channel of inbound messages from the subscribers
	Reg       chan *Cli                // channel to register client to given topics
	UnReg     chan *Cli                // channel to unregister client from given topics when they close connection
	clients   int64                    // number of clients connected to the hub
}

// client struct
type Cli struct {
	Conn       *websocket.Conn // websocket connection
	Registered bool            // is client registered to hub
	Send       chan Payload    // channel of messages to send to the client from the hub
	Topic      string          // topic subscribed to by the client
	tid        string          // trace id
}

// payload to be sent to client and receive from client
type Payload struct {
	Indent int64  `json:"indent"`
	Topic  string `json:"topic"`
}
