package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
)

func WebsocketHandler() fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		client := &Client{conn: conn, send: make(chan []byte, 256)}
		DefaultHub.register <- client
		defer func() {
			DefaultHub.unregister <- client
			if err := conn.Close(); err != nil {
				log.Error().Err(err).Msg("Failed to close connection")
			}
		}()

		go client.writePump()
		client.readPump(DefaultHub)
	})
}
