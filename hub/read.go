package hub

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cms "github.com/tashi21/websocket/commons"
)

// read messages from the client and pass them to the hub
func (c *Cli) ReadPump(ctx *gin.Context) {
	c.tid = cms.Tid(ctx)
	cadr := c.Conn.RemoteAddr().String()

	defer func() {
		log.Printf("[hub.Cli.ReadPump::OK Starting to close connection [tid=%s [ip=%s [clients_registered=%d", c.tid, cadr, H.GetClients())
		H.UnReg <- c   // unregister client from hub
		c.Conn.Close() // close the connection to the client
		log.Printf("[hub.Cli.ReadPump::OK Connection closed [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
	}()

	c.Conn.SetReadLimit(cms.MaxMsgSize)                  // set max message size that can be read from the client
	c.Conn.SetReadDeadline(time.Now().Add(cms.PongWait)) // set read deadline
	c.Conn.SetPongHandler(func(string) error {           // handle pong message from the client
		c.Conn.SetReadDeadline(time.Now().Add(cms.PongWait))
		return nil
	})

	for {
		var msg Payload              // payload to be recieved from the client
		err := c.Conn.ReadJSON(&msg) // read message from the client
		if err != nil {              // if error occured while reading from client, close connection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) { // if error is not due to client closing the connection
				log.Printf("[hub.Cli.ReadPump::KO Unexpected connection close [tid=%s [client=%s [topic=%s [err=%v [clients_registered=%d", c.tid, cadr, c.Topic, err, H.GetClients())
			}
			break
		}

		// if topic is not empty, register the client to the hub
		if msg.Topic != "" {
			if c.Registered { // if client is already registered
				if c.Topic == msg.Topic { // check if client is registered to the same topic
					H.Broadcast <- msg // broadcast the message to the hub
				} else { // if client is registered to a differemt topic, close connection
					c.Write(websocket.CloseMessage, []byte(fmt.Sprintf("client registered to topic '%s' can't register to '%s'", c.Topic, msg.Topic))) // write error message to client
					log.Printf("[hub.Cli.ReadPump::KO Client registered to topic '%s' can't register to '%s' [tid=%s [client=%s [topic=%s [clients_registered=%d", c.Topic, msg.Topic, c.tid, cadr, c.Topic, H.GetClients())
					break
				}
			} else { // if client is not registered
				c.Topic = msg.Topic // set topic of the client
				H.Reg <- c          // register the client to the hub
				H.Broadcast <- msg  // broadcast the payload to the hub
			}
		} else { // if topic is empty, close connection to client
			c.Write(websocket.CloseMessage, []byte("subscribe to at least one topic")) // write error message to client
			log.Printf("[hub.Cli.ReadPump::KO Empty topic [tid=%s [client=%s [topic=%s [clients_registered=%d", c.tid, cadr, c.Topic, H.GetClients())
			break
		}
	}
}
