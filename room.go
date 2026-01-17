package main

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/gofiber/websocket/v2"
)

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

var rooms = make(map[string]*room)

var mu sync.Mutex

func getRoom(name string) *room {

	mu.Lock()
	defer mu.Unlock()

	if r, ok := rooms[name]; ok {
		return r
	}
	r := newRoom()
	rooms[name] = r

	go r.run()
	return r

}

const (
	messageBufferSize = 256
)

func (r *room) Serve(socket *websocket.Conn) {
	roomName := socket.Query("room")
	if roomName == "" {
		_ = socket.WriteMessage(websocket.TextMessage, []byte("room name required"))
		_ = socket.Close()
		return
	}
	realRoom := getRoom(roomName)
	name := socket.Query("name")
	if name == "" {
		name = fmt.Sprintf("USER_%d", rand.Int())
	}

	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    realRoom,
		name:    name,
	}

	realRoom.join <- client
	defer func() { realRoom.leave <- client }()

	go client.write()
	client.read()

}
