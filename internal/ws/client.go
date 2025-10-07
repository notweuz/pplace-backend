package ws

import (
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close connection")
		}
	}()

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	_ = c.conn.Close()
}
