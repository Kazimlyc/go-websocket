package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

func main() {
	var addr = flag.String("addr", ":8080", "Addr of the app")
	flag.Parse()

	app := fiber.New()
	r := newRoom()

	go r.run()

	log.Println("Starting web server on:", *addr)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	app.Get("/chat.html", func(c *fiber.Ctx) error {
		return c.SendFile("chat.html")
	})

	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		r.Serve(conn)
	}))

	if err := app.Listen(*addr); err != nil {
		log.Fatal("Listen and serve:", err)
	}

}
