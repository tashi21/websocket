package hub

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	cms "github.com/tashi21/websocket/commons"
)

// write JSON message from the hub to the client
func (c *Cli) WriteJSON(payload Msg) error {
	c.Conn.SetWriteDeadline(time.Now().Add(cms.WriteWait)) // update write deadline
	return c.Conn.WriteJSON(payload)
}

// write messages from the hub to the client
func (c *Cli) Write(mt int, data []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(cms.WriteWait)) // update write deadline
	return c.Conn.WriteMessage(mt, data)
}

// write messages from the hub to the client
func (c *Cli) WritePump() {
	ticker := time.NewTicker(cms.PingPeriod) // ticker to send ping message to client
	cadr := c.Conn.RemoteAddr().String()

	defer func() { // only invoked when send channel of client is closed
		log.Printf("[hub.Cli.WritePump:: Starting to close connection [client=%s [topic=%s [clients_registered=%d", cadr, c.Topic, H.GetClients())
		ticker.Stop()  // stop ticker
		H.UnReg <- c   // unregister client from hub
		c.Conn.Close() // close the connection
		log.Printf("[hub.Cli.WritePump::OK Connection closed [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
	}()

	for {
		select {
		case msg, ok := <-c.Send: // recieve message from the hub
			if !ok { // if channel is closed, close connection
				c.Write(websocket.CloseMessage, []byte{}) // write empty payload, with close message type
				log.Printf("[hub.Cli.WritePump::OK Channel closed [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
				return
			}

			err := c.WriteJSON(msg)
			if err != nil { // if error occured while writing to client, close connection
				log.Printf("[hub.Cli.WritePump::KO Error writing to client [tid=%s [err=%v [client=%s [topic=%s [clients_registered=%d", c.tid, err, cadr, c.Topic, H.GetClients())
				return
			} else {
				log.Printf("[hub.Cli.WritePump::OK Message sent to client [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
			}

		case <-ticker.C: // if ticker ticks, send ping to client
			err := c.Write(websocket.PingMessage, []byte{}) // write empty payload
			if err != nil {                                 // if error occured while writing to client, close conection
				log.Printf("[hub.Cli.WritePump:: Error pinging client [tid=%s [err=%v [client=%s [topic=%s [clients_registered=%d", c.tid, err, cadr, c.Topic, H.GetClients())
				return
			} else {
				log.Printf("[hub.Cli.WritePump::OK Ping sent to client [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
			}
		}
	}
}
