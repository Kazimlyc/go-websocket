package main

import "github.com/gofiber/websocket/v2"

type room struct {
	clients map[*client]bool
	join    chan *client
	leave   chan *client
	forward chan message
}

type message struct {
	sender *client
	data   []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				if client == msg.sender {
					continue
				}
				client.receive <- msg.data
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

func (r *room) Serve(socket *websocket.Conn) {

	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- client
	defer func() { r.leave <- client }()

	go client.write()
	client.read()

}
