package hub

import "log"

var H = Hub{
	Broadcast: make(chan Msg),
	Reg:       make(chan *Cli),
	UnReg:     make(chan *Cli),
	Topics:    make(map[string]map[*Cli]bool),
	clients:   0,
}

// start up the hub, infinite for loop never breaks
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Reg: // register client to the hub for given topic
			clients := h.Topics[c.Topic] // get clients subscribed to the topic
			if clients == nil {          // if no clients are subscribing to this topic
				clients = make(map[*Cli]bool) // make new map of clients
				h.Topics[c.Topic] = clients   // assign empty map to topic
			}
			h.Topics[c.Topic][c] = true // add client to topic's map of clients
			c.Registered = true         // mark client as registered
			h.AddClient(1)              // increment number of clients connected to hub
			log.Printf("[hub.Hub.Run::OK Registered client [tid=%s [topic=%s [client=%s [clients_registered=%d", c.tid, c.Topic, c.Conn.RemoteAddr(), H.GetClients())

		case c := <-h.UnReg: // unregister client from the hub
			clients := h.Topics[c.Topic] // get clients subscribed to the topic
			if clients != nil {          // if the topic has clients
				if _, ok := clients[c]; ok { // if client exists in map of clients of the topic
					delete(clients, c)   // remove client from map of clients of the topic
					close(c.Send)        // close client's send channel
					h.RemoveClient(1)    // decrement number of clients connected to the hub
					c.Registered = false // mark client as unregistered
					log.Printf("[hub.Hub.Run::OK Unregistered client [tid=%s [topic=%s [client=%s [clients_registered=%d", c.tid, c.Topic, c.Conn.RemoteAddr(), H.GetClients())
					if len(clients) == 0 { // if no clients left
						delete(h.Topics, c.Topic) // remove topic from hub's map of topics
						log.Printf("[hub.Hub.Run::OK Removed topic [tid=%s [topic=%s [clients_registered=%d", c.tid, c.Topic, H.GetClients())
					}
				}
			}

		case m := <-h.Broadcast: // broadcast message to all clients subscribing to same topic
			clients := h.Topics[m.Topic] // get clients subscribed to given message's topic
			for c := range clients {     // for each client subscribed to topic
				cadr := c.Conn.RemoteAddr().String()
				select {
				case c.Send <- m: // send message to client's send channel
					log.Printf("[hub.Hub.Run::OK Broadcasted message [tid=%s [topic=%s [client=%s [clients_registered=%d", c.tid, m.Topic, cadr, H.GetClients())
				default: // if send channel of the client is not receiving
					log.Printf("[hub.Hub.Run::KO Failed to broadcast message [tid=%s [topic=%s [client=%s [clients_registered=%d", c.tid, m.Topic, cadr, H.GetClients())
					h.UnReg <- c // unregister client from hub
				}
			}
		}
	}
}

// add client to hub
func (h *Hub) AddClient(i int64) {
	h.clients += i
}

// remove client from hub
func (h *Hub) RemoveClient(i int64) {
	h.clients -= i
}

// get number of clients connected to hub
func (h Hub) GetClients() int64 {
	return h.clients
}
