package main

import "github.com/gofiber/websocket/v2"

type client struct {
	scoker  *websocket.Conn
	receive chan []byte
	room    *room
}
