package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cms "github.com/tashi21/websocket/commons"
	"github.com/tashi21/websocket/hub"
)

// POST call to broadcast a message into the hub
func SendMsg(ctx *gin.Context) {
	// logging variables
	start := time.Now()
	ip := ctx.ClientIP()
	tid := cms.Tid(ctx)

	// read request body
	var msg hub.Msg
	err := ctx.BindJSON(&msg)
	if err != nil {
		log.Printf("[api.SendMsg::KO Could not bind request body [tid=%s [err=%v [clients_registered=%d", tid, err, hub.H.GetClients())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if msg.Topic != "" { // if topic is not empty
		// send message to hub
		hub.H.Broadcast <- msg
		log.Printf("[api.SendMsg::OK [tid=%s [time_taken=%v [ip=%s [clients_registered=%d", tid, time.Since(start), ip, hub.H.GetClients())
		ctx.Status(http.StatusOK)
	} else { // if topic is empty
		log.Printf("[api.SendMsg::KO Topic is empty [tid=%s [time_taken=%v [ip=%s [clients_registered=%d", tid, time.Since(start), ip, hub.H.GetClients())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing topic"})
	}
}

// reverse proxy for websocket, upgrades connection from GET to WS protocol
func WS(ctx *gin.Context) {
	// logging variables
	start := time.Now()
	ip := ctx.ClientIP()
	tid := cms.Tid(ctx)
	topic := ctx.Param("topic") // topic to subscribe to

	// upgrade HTTP connection to a websocket
	conn, err := cms.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("[api.WS::KO Could not upgrade connection [tid=%s [err=%v [clients_registered=%d", tid, err, hub.H.GetClients())
		return
	}
	log.Printf("[api.WS:: Connection Upgraded [tid=%s [topic=%s [clients_registered=%d", tid, topic, hub.H.GetClients())

	serveWS(ctx, conn, topic) // handle websocket requests
	log.Printf("[api.WS::OK [tid=%s [time_taken=%v [ip=%s [clients_registered=%d", tid, time.Since(start), ip, hub.H.GetClients())
}

// handle websocket requests
func serveWS(ctx *gin.Context, conn *websocket.Conn, topic string) {
	c := &hub.Cli{ // create a new client
		Conn:       conn,
		Registered: false,
		Send:       make(chan hub.Msg, 256),
		Topic:      topic,
	}

	go c.WritePump()          // pump messages from the hub to the client
	go c.ReadPump(ctx, topic) // pump messages from the client to the hub
}
