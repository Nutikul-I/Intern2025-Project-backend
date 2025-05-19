package main

import (
	"payso-internal-api/repository"
	"payso-internal-api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// โหลด .env
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("อ่าน .env ไม่ได้: ", err)
	}
	repository.Init()
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	router.SetupRoutes(app)
	log.Infof("==-- Start Internal Service --== :%s", viper.GetString("SERVER_PORT"))
	log.Fatal(app.Listen(":" + viper.GetString("SERVER_PORT")))
}
