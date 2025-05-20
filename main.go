package main

import (
	"payso-internal-api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	// เริ่มต้นการตั้งค่า log
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("Starting API server...")

	// สร้าง Fiber app
	app := fiber.New(fiber.Config{
		// เพิ่มการตั้งค่าตามที่ต้องการ
		AppName:      "Payso Internal API",
		ErrorHandler: customErrorHandler,
	})

	// เพิ่ม middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// ตั้งค่า routes
	router.SetupRoutes(app)

	// เพิ่ม simple root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Payso Internal API is running")
	})

	// เริ่มการรัน server
	log.Info("Server is running on port 8080")
	log.Fatal(app.Listen(":8080"))
}

// Custom error handler
func customErrorHandler(c *fiber.Ctx, err error) error {
	// ส่ง default error 500 กลับไป
	code := fiber.StatusInternalServerError

	// Override status code if known error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": err.Error(),
		"data":    nil,
	})
}
