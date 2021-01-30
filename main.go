package main

import (
	"log"
	"websockets/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	server := fiber.New()

	server.Use(favicon.New())
	server.Use(logger.New())

	server.Static("/", "./public/index.html")

	util.SetupWS(server)

	server.Get("*", func(c *fiber.Ctx) error {
		return c.Status(404).SendString("sorry bro, 404")
	})

	log.Fatal(server.Listen(":3004"))
}
