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
	var pyld hub.Payload
	err := ctx.BindJSON(&pyld)
	if err != nil {
		log.Printf("[api.SendMsg::KO Could not bind request body [tid=%s [err=%v [clients_registered=%d", tid, err, hub.H.GetClients())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if pyld.Topic != "" { // if topic is not empty
		// send message to hub
		hub.H.Broadcast <- pyld
		log.Printf("[api.SendMsg::OK [tid=%s [time_taken=%v [ip=%s [clients_registered=%d", tid, time.Since(start), ip, hub.H.GetClients())
		ctx.Status(http.StatusOK)
	} else { // if topic is empty
		log.Printf("[api.SendMsg::KO Topic is empty [tid=%s [time_taken=%v [ip=%s [clients_registered=%d", tid, time.Since(start), ip, hub.H.GetClients())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "topic is empty"})
	}
}

func WS(ctx *gin.Context) {
	// logging variables
	start := time.Now()
	ip := ctx.ClientIP()
	tid := cms.Tid(ctx)
	// upgrade HTTP connection to a websocket
	conn, err := cms.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("[api.WS::KO Could not upgrade connection [tid=%s [err=%v [clients_registered=%d", tid, err, hub.H.GetClients())
		return
	}

	serveWS(ctx, conn) // handle websocket requests
	log.Printf("[api.WS::OK [tid=%s [time_taken=%v [ip=%s [clients_registered=%d", tid, time.Since(start), ip, hub.H.GetClients())
}

// handle websocket requests
func serveWS(ctx *gin.Context, conn *websocket.Conn) {
	c := &hub.Cli{ // create a new client
		Conn:       conn,
		Registered: false,
		Send:       make(chan hub.Payload, 256),
		Topic:      "",
	}

	go c.WritePump()   // pump messages from the hub to the client
	go c.ReadPump(ctx) // pump messages from the client to the hub
}
