package router

import (
	"basket-buddy-backend/handler"
	"basket-buddy-backend/middleware"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, client *firestore.Client) {
	setUpWellKnownRoutes(app)
	setUpShareRedirectRoutes(app)

	api := app.Group("/api", logger.New())
	setUpShareRoutes(api, client)
	api.Get("/health", handler.HealthEndpoint())

}

func setUpShareRoutes(api fiber.Router, client *firestore.Client) {
	config := api.Group("/v1/share")
	config.Use(middleware.IsExpired())

	config.Post("/", handler.CreateShareEndpoint(client))
	config.Get("/:ShareCode", handler.FetchShareEndpoint(client))
}

func setUpWellKnownRoutes(app *fiber.App) {
	wellKnown := app.Group("/.well-known")
	wellKnown.Get("/apple-app-site-association", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json; charset=utf-8")

		content, err := os.ReadFile("apple-app-site-association.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading AASA file")
		}

		return c.Send(content)
	})
}

func setUpShareRedirectRoutes(app *fiber.App) {
	app.Get("/share", func(c *fiber.Ctx) error {
		return c.Redirect("https://apps.apple.com/us/app/basketbuddy/id6446040498")
	})
}
