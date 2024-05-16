package vgin

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	webSocketKey     = "VGIN-WEBSOCKET-KEY"
	webSocketConnKey = "VGIN-WEBSOCKET-CONN-KEY"
)

var (
	ErrWebsocketNotSupported      = errors.New("api not support websocket connection")
	ErrWebsocketUpgraderNotExists = errors.New("websocket upgrader not exists")
)

type WebSocketHandler interface {
	Handler
	HandleWebSocket(ctx context.Context, c *Context, conn *websocket.Conn) HandleResponse
}

type WebSocketIsolatedHandler interface {
	WebSocketHandler
	IsolatedHandler
}

type WebSocketInject struct{}

func (ws *WebSocketInject) HandleFunc(ctx context.Context, c *Context) HandleResponse {
	if !websocket.IsWebSocketUpgrade(c.Request) {
		return ErrorRet(http.StatusBadRequest, ErrWebsocketNotSupported, ErrWebsocketNotSupported.Error())
	}
	temp, exists := c.Get(webSocketKey)
	if !exists {
		return ErrorRet(http.StatusBadRequest, ErrWebsocketUpgraderNotExists, ErrWebsocketUpgraderNotExists.Error())
	}
	upgrader := temp.(*websocket.Upgrader)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return ErrorRet(http.StatusInternalServerError, err, "websocket upgrade handshake error")
	}

	c.Set(webSocketConnKey, conn)

	return nil
}

func RangeSendMessage(conn *websocket.Conn, dataGetter func() any) error {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for range tick.C {
		b, err := json.Marshal(dataGetter())
		if err != nil {
			return err
		}
		if err = conn.WriteMessage(websocket.TextMessage, b); err != nil {
			if strings.Contains(err.Error(), "broken pipe") {
				break
			}
			return err
		}
	}

	return nil
}
