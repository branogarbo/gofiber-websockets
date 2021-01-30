package util

import (
	"encoding/json"
	"fmt"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupWS(server *fiber.App) {
	clients := make(map[string]int)

	server.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	////////////////////////////////////

	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {
		fmt.Println(ep.Kws.GetAttribute("UUID"), "connected")
	})
	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		delete(clients, ep.Kws.GetAttribute("UUID"))
		fmt.Println(ep.SocketAttributes["UUID"], "disconnected")
	})
	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		delete(clients, ep.SocketAttributes["UUID"])
		fmt.Println(ep.SocketAttributes["UUID"], "connection closed")
	})
	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		fmt.Println("ikisocket error occurred")
	})
	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		fmt.Println("recieved message from", ep.SocketAttributes["UUID"])

		message := Message{}

		err := json.Unmarshal(ep.Data, &message)
		if err != nil {
			fmt.Println(err)
			return
		}

		if message.Data == "ping" {
			reply, err := json.Marshal(Message{
				Data: "pong",
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			ep.Kws.Emit(reply)
			fmt.Println("sent message to", ep.SocketAttributes["UUID"])
		}

	})

	////////////////////////////////////

	server.Get("/ws", ikisocket.New(func(kws *ikisocket.Websocket) {
		clients[kws.UUID] = len(clients)

		kws.SetAttribute("UUID", kws.UUID)

		kws.Broadcast([]byte(fmt.Sprintf("%s connected", kws.UUID)), true)
	}))
}
