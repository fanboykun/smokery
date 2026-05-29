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

// RegisterWebSocket registers run event streams. The /api/runs/:id/events
// route matches the URL returned by GET /api/runs/{id}. The legacy /ws/runs/:id
// route remains available for older clients.
func RegisterWebSocket(e *echo.Echo, bus port.EventBus) {
	handler := func(c echo.Context) error {
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
	}

	e.GET("/api/runs/:id/events", handler)
	e.GET("/ws/runs/:id", handler)
}
