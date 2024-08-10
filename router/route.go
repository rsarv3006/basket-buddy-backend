package router

import (
	"basket-buddy-backend/ent"
	"basket-buddy-backend/handler"
	"basket-buddy-backend/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, dbClient *ent.Client) {
	setUpWellKnownRoutes(app)
	setUpShareRedirectRoutes(app)

	api := app.Group("/api", logger.New())
	setUpShareRoutes(api, dbClient)
	api.Get("/health", handler.HealthEndpoint())

}

func setUpShareRoutes(api fiber.Router, dbClient *ent.Client) {
	config := api.Group("/v1/share")
	config.Use(middleware.IsExpired())

	config.Post("/", handler.CreateShareEndpoint(dbClient))
	config.Get("/:ShareCode", handler.FetchShareEndpoint(dbClient))
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
