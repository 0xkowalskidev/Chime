package main

import (
	"0xKowalskiDev/Chime/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/html/v2"
)

func main() {
	db, err := db.InitDB("chime.db")
	if err != nil {
		panic(err)
	}
	defer db.DB.Close()

	// Setup Fiber
	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(compress.New()) // Enable brotli compression
	app.Static("/", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		chatrooms, err := db.GetChatrooms()
		if err != nil {
			panic(err) //handle err
		}
		messages, err := db.GetMessages(1)
		if err != nil {
			panic(err)
		}

		return c.Render("index", fiber.Map{"CurrentChatroom": chatrooms[0].Name, "Chatrooms": chatrooms, "Messages": messages})
	})

	log.Fatal(app.Listen(":3000"))
}
