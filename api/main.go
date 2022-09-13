package main

import (
	"fiber-api/configs"
	"fiber-api/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	configs.ConnectDB()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "*",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	routes.UserRoute(app)

	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": os.Getenv("PASSWORD"),
		},
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			if user == "admin" && pass == os.Getenv("PASSWORD") {
				return true
			}
			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{"status": fiber.StatusUnauthorized, "data": "Unauthorized!"})
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}))

	routes.AdminRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello Fiber"})
	})

	app.Listen(":6220")
}
