package main

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/superwhys/venkit/vflags"
	"github.com/superwhys/venkit/vgin"
)

type TestWebSocketHandler struct {
	*vgin.WebSocketInject
	UserId int `form:"user_id"`
}

func (h *TestWebSocketHandler) InitHandler() vgin.IsolatedHandler {
	return new(TestWebSocketHandler)
}

func (h *TestWebSocketHandler) HandleWebSocket(ctx context.Context, c *vgin.Context, conn *websocket.Conn) vgin.HandleResponse {
	if err := c.ShouldBind(h); err != nil {
		return vgin.ErrorRet(400, err, "err")
	}
	vgin.RangeSendMessage(conn, func() any {
		return vgin.Data{"data": h.UserId}
	})
	return nil
}

func main() {
	vflags.Parse()
	engine := vgin.New()
	engine.GET("/wx", &TestWebSocketHandler{})

	engine.Run(":2888")
}
