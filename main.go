package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/html/v2"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Init SQLite
	db, err := sql.Open("sqlite3", "./chime.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create table if it does not exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS chatrooms(
                id INTEGER
        )
        `)
	if err != nil {
		panic(err)
	}

	// Setup Fiber
	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(compress.New()) // Enable brotli compression
	app.Static("/", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	log.Fatal(app.Listen(":3000"))
}
