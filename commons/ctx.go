package commons

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	TraceKey string = "req-trace-id"
)

func NewTid() int64 {
	return int64(uuid.New().ID())
}

func Tid(ctx *gin.Context) string {
	tid := ctx.GetString(TraceKey)
	if len(tid) == 0 {
		tid = strconv.FormatInt(NewTid(), 10)
		ctx.Set(TraceKey, tid)
	}
	return tid
}
