package main

import (
	"payso-internal-api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	log.Info("==-- Start Internal Service --==")
	router.SetupRoutes(app)

	app.Listen(":" + viper.GetString("SERVER_PORT"))

}
