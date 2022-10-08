package commons

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pelletier/go-toml"
	glcms "gitlab.com/raaho/golangcommons"
)

var (
	HttpPortInt uint
	HttpHostInt string
	HttpPortExt uint
	HttpHostExt string
	SrvConf     *toml.Tree
	Upgrader    = websocket.Upgrader{ // use default options
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	confEnv    = "SERVICE_CONF"
	dfConf     = "conf/service.toml"
	WriteWait  = 10 * time.Second    // time allowed to write a message to the client
	PongWait   = 60 * time.Second    // time allowed to read the next pong message from the client
	PingPeriod = (PongWait * 9) / 10 // send pings to client with this period
	MaxMsgSize = 512                 // max message size (in bytes) allowed from the client
)

func init() {
	SrvConf = glcms.LoadCfg(confEnv, dfConf)

	// init attrs
	HttpHostInt = SrvConf.Get("http.int.host").(string)
	HttpPortInt = uint(SrvConf.Get("http.int.port").(int64)) // converting to uint as it interface supports only int64

	HttpHostExt = SrvConf.Get("http.ext.host").(string)
	HttpPortExt = uint(SrvConf.Get("http.ext.port").(int64)) // converting to uint as it interface supports only int64
}
