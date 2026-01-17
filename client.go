package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

type client struct {
	socket  *websocket.Conn
	receive chan []byte
	room    *room
	name    string
}

func (c *client) read() {
	for {

		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}

		//c.room.forward <- message{sender: c, data: msg}
		outgoing := map[string]string{
			"name":    c.name,
			"message": string(msg),
		}
		jsMessage, err := json.Marshal(outgoing)
		if err != nil {
			fmt.Println("encoding failed", err)
			continue
		}

		c.room.forward <- message{sender: c, data: jsMessage}
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
