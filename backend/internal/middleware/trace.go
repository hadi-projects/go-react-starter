package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const TraceKey = "middleware_trace"

func AddToTrace(ctx *gin.Context, name string) {
	trace, exists := ctx.Get(TraceKey)
	var traceList []string
	if exists {
		traceList = trace.([]string)
	}
	traceList = append(traceList, name)
	ctx.Set(TraceKey, traceList)
}

func GetTraceString(ctx *gin.Context) string {
	trace, exists := ctx.Get(TraceKey)
	if !exists {
		return ""
	}
	traceList := trace.([]string)
	return strings.Join(traceList, " -> ")
}
