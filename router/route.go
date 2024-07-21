package router

import (
	"basket-buddy-backend/ent"
	"basket-buddy-backend/handler"
	"basket-buddy-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, dbClient *ent.Client) {
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
