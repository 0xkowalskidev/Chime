package main

import (
	DB "0xKowalskiDev/Chime/db"
	"bytes"
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/html/v2"
)

type Client struct {
	Conn       *websocket.Conn
	ChatroomID int
}

var (
	clients   = make(map[*Client]bool)
	clientsMu sync.Mutex
)

func main() {
	db, err := DB.InitDB("chime.db")
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

	// Index
	app.Get("/", func(c *fiber.Ctx) error {
		chatrooms, err := db.GetChatrooms()
		if err != nil {
			panic(err) //handle err
		}
		messages, err := db.GetMessages(1)
		if err != nil {
			panic(err)
		}

		return c.Render("index", fiber.Map{"CurrentChatroomID": chatrooms[0].ID, "CurrentChatroom": chatrooms[0].Name, "Chatrooms": chatrooms, "Messages": messages})
	})

	// Websockets
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:chatroom", websocket.New(func(c *websocket.Conn) {
		chatroomID, _ := strconv.Atoi(c.Params("chatroom")) // assuming cant error, shouldent
		client := &Client{Conn: c, ChatroomID: chatroomID}

		// Register client
		clientsMu.Lock()
		clients[client] = true
		clientsMu.Unlock()

		defer func() {
			clientsMu.Lock()
			delete(clients, client)
			clientsMu.Unlock()
			c.Close()
		}()

		// Handle incoming messages
		for {
			var msg DB.Message
			err := c.ReadJSON(&msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway) {
					break
				}
				log.Println("Read error:", err)
				break
			}

			// Save to DB
			message, err := db.CreateMessage(chatroomID, "user", msg.Content)
			if err != nil {
				log.Println("DB error:", err)
				continue
			}

			// Broadcast to all clients in the same chatroom
			var buf bytes.Buffer
			if err := engine.Render(&buf, "message", message); err != nil {
				log.Println("Render error:", err)
				continue
			}

			messageHTML := []byte(`<div hx-swap-oob="beforeend:#chat-messages">` + buf.String() + `</div>`)

			broadcastMessage(chatroomID, messageHTML)
		}
	}))

	log.Fatal(app.Listen(":3000"))
}

func broadcastMessage(chatroomID int, messageHTML []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		if client.ChatroomID == chatroomID {
			if err := client.Conn.WriteMessage(websocket.TextMessage, messageHTML); err != nil {
				log.Println("Write error:", err)
				client.Conn.Close()
				delete(clients, client)
			}
		}
	}
}
