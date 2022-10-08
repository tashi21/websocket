package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	cms "github.com/tashi21/websocket/commons"
	"github.com/tashi21/websocket/expapi/api"
	"github.com/tashi21/websocket/hub"
	"gitlab.com/raaho/apicommon"
	"gitlab.com/raaho/apicommongin"
)

func main() {
	// defer dbutil.Close() // initialize DB before any request
	go hub.H.Run() // initializ the hub
	initInt()      // mapping for internal apis
	initExt()      // mapping for external apis
}

// REST External::APIs for Apps/Web Consoles
func initExt() {
	rExt := apicommongin.NewRouter()
	pExt := apicommon.CliPortWT(cms.HttpPortExt, "ext")
	router(&rExt.RouterGroup)

	rExt.Run(fmt.Sprintf("%v:%v", cms.HttpHostExt, pExt))
}

// REST Internal::APIs for Services
func initInt() {
	rInt := apicommongin.NewRouter()
	pInt := apicommon.CliPortWT(cms.HttpPortInt, "int")
	router(&rInt.RouterGroup)
	go rInt.Run(fmt.Sprintf("%v:%v", cms.HttpHostInt, pInt))
}

// add routes to given gin router
func router(r *gin.RouterGroup) {
	r.POST("/message", api.SendMsg)
	r.GET("/ws", api.WS)
}
