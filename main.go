package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Data string `json:"data"`
}

func main() {
	clients := make(map[string]string)

	server := fiber.New()

	server.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	////////////////////////////////////

	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {
		fmt.Println(ep.Kws.GetAttribute("user_id"), "connected")
	})

	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		delete(clients, ep.Kws.GetAttribute("user_id"))
		fmt.Println(ep.SocketAttributes["user_id"], "disconnected")
	})

	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		delete(clients, ep.SocketAttributes["user_id"])
		fmt.Println(ep.SocketAttributes["user_id"], "connection closed")
	})

	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		fmt.Println("ikisocket error occurred")
	})
	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		fmt.Println(ep.SocketAttributes["user_id"], "connection closed")

		message := Message{}

		err := json.Unmarshal(ep.Data, &message)
		if err != nil {
			fmt.Println(err)
			return
		}

		if message.Data == "ping" {
			reply, _ := json.Marshal(Message{
				Data: "pong",
			})

			ep.Kws.Emit(reply)
		}

	})

	////////////////////////////////////

	server.Get("/ws/:id", ikisocket.New(func(kws *ikisocket.Websocket) {
		userId := kws.Params("id")

		clients[userId] = kws.UUID

		kws.SetAttribute("user_id", userId)

		kws.Broadcast([]byte(fmt.Sprintf("%s connected", userId)), true)
	}))

	log.Fatal(server.Listen(":3004"))
}
