package hub

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cms "github.com/tashi21/websocket/commons"
)

// read messages from the client and pass them to the hub
func (c *Cli) ReadPump(ctx *gin.Context, topic string) {
	c.tid = cms.Tid(ctx)
	cadr := c.Conn.RemoteAddr().String()

	defer func() {
		log.Printf("[hub.Cli.ReadPump::OK Starting to close connection [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
		H.UnReg <- c   // unregister client from hub
		c.Conn.Close() // close the connection to the client
		log.Printf("[hub.Cli.ReadPump::OK Connection closed [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
	}()

	c.Conn.SetReadLimit(cms.MaxMsgSize)                  // set max message size that can be read from the client
	c.Conn.SetReadDeadline(time.Now().Add(cms.PongWait)) // set read deadline
	c.Conn.SetPongHandler(func(string) error {           // handle pong message from the client
		c.Conn.SetReadDeadline(time.Now().Add(cms.PongWait)) // update read deadline
		return nil
	})

	H.Reg <- c // register the client to the hub

	for {
		var msg Msg                  // payload to be recieved from the client
		err := c.Conn.ReadJSON(&msg) // read message from the client
		if err != nil {              // if error occured while reading from client, close connection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) { // if error is not due to client closing the connection
				log.Printf("[hub.Cli.ReadPump::KO Unexpected connection close [tid=%s [client=%s [topic=%s [err=%v [clients_registered=%d", c.tid, cadr, c.Topic, err, H.GetClients())
			}
			break
		}

		H.Broadcast <- msg // broadcast the message to the hub
	}
}
