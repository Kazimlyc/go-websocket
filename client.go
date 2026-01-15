package main

import "github.com/gofiber/websocket/v2"

type client struct {
	scoker  *websocket.Conn
	receive chan []byte
	room    *room
}

func (c *client) read() {
	for {

		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
