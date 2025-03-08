package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(compress.New())

	app.Static("/", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	log.Fatal(app.Listen(":3000"))
}
