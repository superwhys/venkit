package middlewares

import (
	"log"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	log2 "github.com/superwhys/venkit/lg/v2/log"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
	
	logMsg = "| %s %v %s | %13v | %15s | %s %4v %s| %#v"
)

func statusCodeColor(code int) string {
	
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func methodColor(method string) string {
	
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	newLogger := log2.New(
		log2.WithCalldepth(5),
		log2.WithDebugFlag(log.LstdFlags|log.LUTC),
		log2.WithErrorFlag(log.LstdFlags|log.LUTC),
	)
	return func(c *gin.Context) {
		start := time.Now()
		
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		
		c.Next()
		gin.Logger()
		
		clientIp := c.ClientIP()
		method := c.Request.Method
		if raw != "" {
			path = path + "?" + raw
		}
		statusCode := c.Writer.Status()
		spendTime := time.Since(start)
		
		var logFunc func(msg string, v ...any)
		switch {
		case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
			logFunc = newLogger.Infof
		case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest:
			logFunc = newLogger.Warnf
		case statusCode >= http.StatusBadRequest && statusCode <= http.StatusNetworkAuthenticationRequired:
			logFunc = newLogger.Errorf
		default:
			logFunc = newLogger.Errorf
		}
		
		logFunc(
			logMsg,
			statusCodeColor(statusCode), statusCode, reset,
			spendTime,
			clientIp,
			methodColor(method), method, reset,
			path,
		)
	}
}
