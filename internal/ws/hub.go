package ws

import (
	"encoding/json"
	"pplace_backend/internal/model/dto/ws"
	"time"

	"pplace_backend/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog/log"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

var DefaultHub *Hub

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

func Start() {
	DefaultHub = NewHub()
	go DefaultHub.Run()
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to remove client")
			return
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
	err := c.conn.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close connection")
		return
	}
}

func WebsocketHandler() fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		client := &Client{conn: conn, send: make(chan []byte, 256)}
		DefaultHub.register <- client
		defer func() {
			DefaultHub.unregister <- client
			err := conn.Close()
			if err != nil {
				log.Error().Err(err).Msg("Failed to close connection")
				return
			}
		}()

		go client.writePump()
		client.readPump(DefaultHub)
	})
}

func BroadcastPixel(action string, pixel *model.Pixel) {
	if DefaultHub == nil {
		return
	}
	ev := ws.NewPixelEventDto(action, pixel.Color, pixel.User.Username, pixel.ID, pixel.X, pixel.Y, pixel.UserID, time.Now())
	data, err := json.Marshal(ev)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal event")
		return
	}
	select {
	case DefaultHub.broadcast <- data:
		log.Debug().Str("action", action).Msg("Broadcasted pixel")
	default:
	}
}

func BroadcastPixelDelete(id, x, y uint) {
	if DefaultHub == nil {
		return
	}
	ev := ws.NewPixelDeleteEventDto("delete", id, x, y, time.Now())
	data, err := json.Marshal(ev)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal delete event")
		return
	}
	select {
	case DefaultHub.broadcast <- data:
		log.Debug().Str("action", ev.Action).Msg("Broadcasted pixel")
	default:
	}
}
