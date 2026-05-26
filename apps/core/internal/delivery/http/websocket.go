package http

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/fanboykun/smokery/apps/core/internal/port"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// RegisterWebSocket registers /ws/runs/:id which streams run events via the EventBus.
// This is a raw Echo handler (not a huma operation) since WebSockets aren't part of OpenAPI.
func RegisterWebSocket(e *echo.Echo, bus port.EventBus) {
	e.GET("/ws/runs/:id", func(c echo.Context) error {
		runID := c.Param("id")
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		ch := bus.Subscribe(runID)
		defer bus.Unsubscribe(runID, ch)

		for event := range ch {
			if err := ws.WriteJSON(event); err != nil {
				break
			}
			if event.Type == port.EventRunFinished {
				break
			}
		}
		return nil
	})
}
